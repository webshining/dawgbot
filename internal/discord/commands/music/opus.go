package music

import (
	"bufio"
	"encoding/binary"
	"io"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"layeh.com/gopus"
)

func (m *Music) playFile(vc *discordgo.VoiceConnection, filename string, stop chan bool) bool {
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
