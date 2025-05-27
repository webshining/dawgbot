package music

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/bwmarrin/dgvoice"
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
			go c.playAudio(s, i, attachment.URL)
			return
		case "youtubeurl":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ссылка получена, обрабатываю...",
				},
			})
			go c.playAudioYoutube(s, i, opt.Value.(string))
			return
		}
	}
}

func (c *Music) joinVoice(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.VoiceConnection, error) {
	vs, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
	if err != nil {
		c.logger.Error("error getting voice state", zap.Error(err))
		return nil, err
	}
	vc, ok := s.VoiceConnections[vs.GuildID]
	if !ok {
		vc, err = s.ChannelVoiceJoin(i.GuildID, vs.ChannelID, false, false)
		if err != nil {
			c.logger.Error("error joining voice channel", zap.Error(err))
			return nil, err
		}
	}
	return vc, nil
}

func (c *Music) playAudioYoutube(s *discordgo.Session, i *discordgo.InteractionCreate, url string) {
	if cancel, exists := c.playbackCancel[i.GuildID]; exists {
		cancel()
	}

	file, err := c.downloadTempFileYoutube(url)
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.playbackCancel[i.GuildID] = cancel

	stop := make(chan bool)
	go func() {
		<-ctx.Done()
		stop <- true
		cancel()
		delete(c.playbackCancel, i.GuildID)
		os.Remove(file)
	}()

	vc, _ := c.joinVoice(s, i)
	go dgvoice.PlayAudioFile(vc, file, stop)
}

func (c *Music) playAudio(s *discordgo.Session, i *discordgo.InteractionCreate, url string) {
	if cancel, exists := c.playbackCancel[i.GuildID]; exists {
		cancel()
	}

	file, err := c.downloadTempFile(url)
	if err != nil {
		return
	}
	if !c.checkIsAudio(file) {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.playbackCancel[i.GuildID] = cancel

	stop := make(chan bool)
	go func() {
		<-ctx.Done()
		stop <- true
		cancel()
		delete(c.playbackCancel, i.GuildID)
		os.Remove(file)
	}()

	vc, _ := c.joinVoice(s, i)
	go dgvoice.PlayAudioFile(vc, file, stop)
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

func (c *Music) downloadTempFileYoutube(url string) (string, error) {
	tempFile, err := os.CreateTemp("", "discord_upload-*.tmp")
	if err != nil {
		return "", err
	}
	filename := tempFile.Name()
	tempFile.Close()
	os.Remove(filename)

	exec.Command("yt-dlp", "-f", "bestaudio", "-o", filename, url).Run()

	return filename, nil
}

func (c *Music) downloadTempFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tempFile, err := os.CreateTemp("", "discord_upload-*.tmp")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
