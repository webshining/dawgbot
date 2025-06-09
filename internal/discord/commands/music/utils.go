package music

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"layeh.com/gopus"
)

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

func (c *Music) downloadTempFile(url string) *string {
	tempFile, _ := os.CreateTemp("", "discord_upload-*.tmp")
	filename := tempFile.Name()
	defer func() {
		tempFile.Close()
		os.Remove(filename)
	}()

	println(url)
	if strings.Contains(url, "youtube") {
		tempFile.Close()
		os.Remove(filename)

		cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-o", filename, url)
		cmd.Run()
	} else if strings.Contains(url, "discord") {
		resp, err := http.Get(url)
		if err != nil {
			c.logger.Error("failed to download file", zap.Error(err))
			return nil
		}
		defer resp.Body.Close()

		io.Copy(tempFile, resp.Body)

		if !c.checkIsAudio(filename) {
			return nil
		}
	}

	client, _ := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("admin_user", "admin_pass", ""),
		Secure: false,
		Region: "us-east-1",
	})
	generatedName := uuid.New().String()
	client.FPutObject(context.Background(), "dawgbot", generatedName, filename, minio.PutObjectOptions{})

	result := fmt.Sprintf("http://localhost:9000/dawgbot/%s", generatedName)
	return &result
}

func (m *Music) PlayFile(vc *discordgo.VoiceConnection, filename string, stop <-chan struct{}) bool {
	ffmpegcmd := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegout, err := ffmpegcmd.StdoutPipe()
	if err != nil {
		m.logger.Error("StdoutPipe Error", zap.Error(err))
		return true
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	if err := ffmpegcmd.Start(); err != nil {
		m.logger.Error("FFmpeg Start Error", zap.Error(err))
		return true
	}
	defer ffmpegcmd.Process.Kill()

	end := true

	go func() {
		<-stop
		end = false
		ffmpegcmd.Process.Kill()
	}()

	vc.Speaking(true)
	defer vc.Speaking(false)

	send := make(chan []int16, 10)
	defer close(send)

	go func() {
		m.sendPCM(vc, send)
	}()

	for {
		audiobuf := make([]int16, 960*2)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return end
		}
		if err != nil {
			m.logger.Error("error reading from ffmpeg stdout", zap.Error(err))
			return true
		}

		select {
		case send <- audiobuf:
		case <-stop:
			return false
		}
	}

}

func (m *Music) sendPCM(vc *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	opusEncoder, err := gopus.NewEncoder(48000, 2, gopus.Audio)
	if err != nil {
		m.logger.Error("Failed to create Opus encoder", zap.Error(err))
		return
	}

	for {
		recv, ok := <-pcm
		if !ok {
			return
		}

		opus, err := opusEncoder.Encode(recv, 960, 960*4)
		if err != nil {
			m.logger.Error("Failed to encode PCM to Opus", zap.Error(err))
			return
		}

		if !vc.Ready || vc.OpusSend == nil {
			return
		}

		vc.OpusSend <- opus
	}
}
