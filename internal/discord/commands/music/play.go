package music

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func (c *Music) Play(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, opt := range i.ApplicationCommandData().Options {
		switch opt.Name {
		case "file":
			attachmentID := i.ApplicationCommandData().Options[0].Value.(string)
			attachment := i.ApplicationCommandData().Resolved.Attachments[attachmentID]
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Файл получен, обрабатываю...",
				},
			})
			go c.playAudio(i, attachment.URL)
			return
		case "youtubeurl":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ссылка получена, обрабатываю...",
				},
			})
			go c.playAudioYoutube(i, opt.Value.(string))
			return
		}
	}
}

func (c *Music) joinVoice(i *discordgo.InteractionCreate) *discordgo.VoiceConnection {
	vs, err := c.session.State.VoiceState(i.GuildID, i.Member.User.ID)
	if err != nil {
		c.logger.Error("error getting voice state", zap.Error(err))
		return nil
	}

	vc, ok := c.session.VoiceConnections[i.GuildID]
	if !ok {
		vc, err = c.session.ChannelVoiceJoin(i.GuildID, vs.ChannelID, false, false)
		if err != nil {
			c.logger.Error("error joining voice channel", zap.Error(err))
			return nil
		}
	}

	return vc
}

func (c *Music) scheduleAutoDisconnect(guildID string) {
	if timer, exists := c.autoDisconnectTimers[guildID]; exists {
		timer.Stop()
	}

	c.autoDisconnectTimers[guildID] = time.AfterFunc(10*time.Second, func() {
		vc, ok := c.session.VoiceConnections[guildID]
		if ok {
			vc.Disconnect()
		}
		delete(c.autoDisconnectTimers, guildID)
	})
}

func (c *Music) playAudioYoutube(i *discordgo.InteractionCreate, url string) {
	// cancel previous playback if exists
	if timer, exists := c.autoDisconnectTimers[i.GuildID]; exists {
		timer.Stop()
		delete(c.autoDisconnectTimers, i.GuildID)
	}
	if cancel, exists := c.playbackCancel[i.GuildID]; exists {
		cancel()
	}

	file := c.downloadTempFileYoutube(url)
	if file == nil {
		return
	}

	vc := c.joinVoice(i)
	if vc == nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.playbackCancel[i.GuildID] = cancel

	// call when cancel() is called
	stop := make(chan bool)
	go func() {
		<-ctx.Done()
		stop <- true
		delete(c.playbackCancel, i.GuildID)
		os.Remove(*file)
	}()

	go func() {
		end := c.playFile(vc, *file, stop)
		if end {
			c.scheduleAutoDisconnect(i.GuildID)
		}
		cancel()
	}()
}

func (c *Music) playAudio(i *discordgo.InteractionCreate, url string) {
	// cancel previous playback if exists
	if timer, exists := c.autoDisconnectTimers[i.GuildID]; exists {
		timer.Stop()
		delete(c.autoDisconnectTimers, i.GuildID)
	}
	if cancel, exists := c.playbackCancel[i.GuildID]; exists {
		cancel()
	}

	file := c.downloadTempFile(url)
	if file == nil {
		return
	}
	if !c.checkIsAudio(*file) {
		return
	}

	vc := c.joinVoice(i)
	if vc == nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.playbackCancel[i.GuildID] = cancel

	// call when cancel() is called
	stop := make(chan bool)
	go func() {
		<-ctx.Done()
		stop <- true
		delete(c.playbackCancel, i.GuildID)
		os.Remove(*file)
	}()

	go func() {
		end := c.playFile(vc, *file, stop)
		if end {
			c.scheduleAutoDisconnect(i.GuildID)
		}
		cancel()
	}()
}

func (c *Music) checkIsAudio(filename string) bool {
	type FFProbeOutput struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
		} `json:"streams"`
	}

	cmd, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "stream=codec_type", "-of", "json", filename).Output()
	if err != nil {
		return false
	}
	var result FFProbeOutput
	if err := json.Unmarshal(cmd, &result); err != nil {
		return false
	}
	for _, stream := range result.Streams {
		if stream.CodecType == "audio" {
			return true
		}
	}
	return false
}

func (c *Music) downloadTempFileYoutube(url string) *string {
	tempFile, err := os.CreateTemp("", "discord_upload-*.tmp")
	if err != nil {
		c.logger.Error("failed to create temp file", zap.Error(err))
		return nil
	}
	filename := tempFile.Name()
	tempFile.Close()
	os.Remove(filename)

	exec.Command("yt-dlp", "-f", "bestaudio", "-o", filename, url).Run()

	return &filename
}

func (c *Music) downloadTempFile(url string) *string {
	resp, err := http.Get(url)
	if err != nil {
		c.logger.Error("failed to download file", zap.Error(err))
		return nil
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp("", "discord_upload-*.tmp")
	if err != nil {
		c.logger.Error("failed to create temp file", zap.Error(err))
		return nil
	}
	filename := tempFile.Name()
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		c.logger.Error("failed to copy response body to temp file", zap.Error(err))
		return nil
	}

	return &filename
}
