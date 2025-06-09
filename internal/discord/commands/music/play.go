package music

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/webshining/internal/common/database"
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
			go c.playHandler(i, attachment.URL)
			return
		case "youtubeurl":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ссылка получена, обрабатываю...",
				},
			})
			go c.playHandler(i, opt.Value.(string))
			return
		}
	}
}

func (c *Music) playHandler(i *discordgo.InteractionCreate, url string) {
	file := c.downloadTempFile(url)
	if file == nil {
		return
	}

	c.db.Create(&database.Playlist{GuildID: i.GuildID, FileUrl: *file})

	if _, ok := c.playbackCancel[i.GuildID]; ok {
		return
	}

	// stop autodisconnect timer
	if timer, exists := c.autoDisconnectTimers[i.GuildID]; exists {
		timer.Stop()
		delete(c.autoDisconnectTimers, i.GuildID)
	}

	vc := c.joinVoice(i)
	if vc == nil {
		return
	}

	go func() {
		end := c.playAudio(vc, *file, vc.GuildID)
		if end {
			c.scheduleAutoDisconnect(i.GuildID)
		}
	}()
}

func (c *Music) playAudio(vc *discordgo.VoiceConnection, file string, guildID string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	c.playbackCancel[guildID] = cancel

	end := c.PlayFile(vc, file, ctx.Done())

	var trackToDelete *database.Playlist
	c.db.Take(trackToDelete, "guild_id = ?", vc.GuildID)
	if trackToDelete != nil {
		c.db.Unscoped().Delete(&trackToDelete)
	}

	var nextTrack *database.Playlist
	c.db.First(nextTrack, "guild_id = ?", vc.GuildID)
	if nextTrack == nil {
		return end
	}

	return c.playAudio(vc, nextTrack.FileUrl, guildID)
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
