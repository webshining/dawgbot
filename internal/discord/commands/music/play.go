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

	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "file":
			attachmentID, _ := i.ApplicationCommandData().Options[0].Value.(string)
			attachment := i.ApplicationCommandData().Resolved.Attachments[attachmentID]
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Файл получен, обрабатываю...",
				},
			})
			c.playAudio(vc, i.GuildID, attachment.URL)
			return
		case "youtubeurl":
			url, _ := opt.Value.(string)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ссылка получена , обрабатываю...",
				},
			})
			c.playAudioYoutube(vc, i.GuildID, url)
			return
		}
	}

}

func (c *Music) playAudioYoutube(vc *discordgo.VoiceConnection, guildID string, url string) {
	if cancel, exists := c.playbackCancel[guildID]; exists {
		cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.playbackCancel[guildID] = cancel
	defer func() {
		cancel()
		delete(c.playbackCancel, guildID)
	}()

	ytdcmd := exec.CommandContext(ctx, "yt-dlp", "-f", "bestaudio", "-o", "-", url)
	ytdOut, _ := ytdcmd.StdoutPipe()

	ffmpegcmd := exec.CommandContext(ctx, "ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegcmd.Stdin = ytdOut
	ffmpegout, _ := ffmpegcmd.StdoutPipe()

	ytdcmd.Start()
	ffmpegcmd.Start()

	defer func() {
		ytdcmd.Process.Kill()
		ffmpegcmd.Process.Kill()
	}()

	// conver pcm to opus
	audiobuf := bufio.NewReaderSize(ffmpegout, 65536)
	opusEncoder, _ := gopus.NewEncoder(48000, 2, gopus.Audio)

	type opusPacket []byte
	type opusBlock []opusPacket

	bufferReady := make(chan struct{}, 1)
	const bufferSize = 100
	audioBuffer := make(chan opusBlock, 10)

	go func() {
		defer close(audioBuffer)
		var block opusBlock

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			raw := make([]int16, 960*2)
			if err := binary.Read(audiobuf, binary.LittleEndian, raw); err != nil {
				if len(block) > 0 {
					audioBuffer <- block
				}
				return
			}

			packet, err := opusEncoder.Encode(raw, 960, 960*4)
			if err != nil {
				c.logger.Error("opus encode error", zap.Error(err))
				if len(block) > 0 {
					audioBuffer <- block
				}
				return
			}

			block = append(block, packet)

			if len(block) == bufferSize {
				select {
				case audioBuffer <- block:
				case <-ctx.Done():
					return
				}
				block = nil
				select {
				case bufferReady <- struct{}{}:
				default:
				}
			}
		}
	}()

	select {
	case <-bufferReady:
	case <-ctx.Done():
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case block, ok := <-audioBuffer:
			if !ok {
				return
			}
			for _, packet := range block {
				select {
				case vc.OpusSend <- packet:
				case <-ctx.Done():
					return
				}
			}
		}
	}
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

	// convert audio to pcm format
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", filename, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegout, err := cmd.StdoutPipe()
	if err != nil {
		c.logger.Error("ffmpeg stdout pipe error", zap.Error(err))
		return
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 65536)

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

	// create a channel to send opus packets
	type opusPacket []byte
	opusChan := make(chan opusPacket, 200)
	go func() {
		for {
			raw := make([]int16, 960*2)
			if err := binary.Read(ffmpegbuf, binary.LittleEndian, raw); err != nil {
				if err != io.EOF {
					c.logger.Error("binary read error", zap.Error(err))
				}
				close(opusChan)
				return
			}

			encoder, err := opusEncoder.Encode(raw, 960, 960*4)
			if err != nil {
				c.logger.Error("opus encode error", zap.Error(err))
				close(opusChan)
				return
			}

			opusChan <- encoder
		}
	}()

	// prebuffering
	prebuffer := 10
	for i := 0; i < prebuffer; i++ {
		_, ok := <-opusChan
		if !ok {
			return
		}
	}

	// start streaming
	for {
		select {
		case <-ctx.Done():
			return
		case packet, ok := <-opusChan:
			if !ok {
				return
			}
			vc.OpusSend <- packet
		}
	}
}
