package music

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"layeh.com/gopus"
)

func (c *Music) Play(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "play" {
		return
	}

	// Получаем ID вложения из опций
	attachmentID, _ := i.ApplicationCommandData().Options[0].Value.(string)
	attachment, _ := i.ApplicationCommandData().Resolved.Attachments[attachmentID]

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Файл получен, обрабатываю...",
		},
	})

	vs, err := c.session.State.VoiceState(i.GuildID, i.Member.User.ID)
	if err != nil {
		c.logger.Error("error getting voice state", zap.Error(err))
		return
	}

	vc, ok := c.session.VoiceConnections[vs.GuildID]
	if !ok {
		vc, err = c.session.ChannelVoiceJoin(i.GuildID, vs.ChannelID, false, false)
		if err != nil {
			c.logger.Error("error joining voice channel", zap.Error(err))
			return
		}
	}

	c.playAudio(vc, i.GuildID, attachment.URL)
}

func (c *Music) playAudio(vc *discordgo.VoiceConnection, guildID string, filename string) {
	if cancel, exists := c.playbackCancel[guildID]; exists {
		cancel()
	}

	// create a new context for the playback
	ctx, cancel := context.WithCancel(context.Background())
	c.playbackCancel[guildID] = cancel
	defer func() {
		cancel()
		delete(c.playbackCancel, guildID)
	}()

	// get file data in ffmpeg format
	cmd := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegout, err := cmd.StdoutPipe()
	if err != nil {
		c.logger.Error("error creating ffmpeg stdout pipe", zap.Error(err))
		return
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	// start ffmpeg
	if err := cmd.Start(); err != nil {
		c.logger.Error("error starting ffmpeg", zap.Error(err))
		return
	}
	defer cmd.Process.Kill()

	vc.Speaking(true)
	defer vc.Speaking(false)

	// create opus encoder
	opusEncoder, err := gopus.NewEncoder(48000, 2, gopus.Audio)
	if err != nil {
		c.logger.Error("error creating opus encoder", zap.Error(err))
		return
	}

	// read from ffmpeg stdout and send to discord
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Audio playback cancelled", zap.String("guild", guildID))
			return
		default:
			audiobuf := make([]int16, 960*2)
			err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return
			}
			if err != nil {
				c.logger.Error("error reading from ffmpeg stdout", zap.Error(err))
				return
			}

			opus, err := opusEncoder.Encode(audiobuf, 960, 960*4)
			if err != nil {
				c.logger.Error("error encoding audio", zap.Error(err))
				return
			}

			if !vc.Ready || vc.OpusSend == nil {
				c.logger.Error("voice connection not ready")
				return
			}
			vc.OpusSend <- opus
		}
	}
}
