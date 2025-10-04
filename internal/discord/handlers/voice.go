package handlers

import (
	"bot/internal/database"
	"fmt"
	"html"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/bwmarrin/discordgo"
)

func (h *handlers) VoiceJoinHandler(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.ChannelID == "" {
		return
	}
	if vs.BeforeUpdate != nil && vs.BeforeUpdate.ChannelID == vs.ChannelID {
		return
	}

	user, err := s.User(vs.UserID)
	if err != nil {
		return
	}
	channel, err := s.State.Channel(vs.ChannelID)
	if err != nil {
		channel, err = s.Channel(vs.ChannelID)
		if err != nil {
			return
		}
	}
	guild, err := s.State.Guild(vs.GuildID)
	if err != nil {
		guild, err = s.Guild(vs.GuildID)
		if err != nil {
			return
		}
	}

	var channelDB *database.Channel
	h.app.DB.Preload("Users").First(&channelDB, vs.ChannelID)

	for _, userDB := range channelDB.Users {
		text := fmt.Sprintf("<code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code>",
			html.EscapeString(guild.Name),
			html.EscapeString(channel.Name),
			html.EscapeString(user.DisplayName()),
		)
		if userDB.LastGuildID != guild.ID {
			userDB.LastGuildID = guild.ID
			h.app.DB.Save(&userDB)
			h.app.TelegramBot.SendPhoto(userDB.ID, gotgbot.InputFileByURL(guild.IconURL("1024")), &gotgbot.SendPhotoOpts{
				Caption:   text,
				ParseMode: "HTML",
			})
		} else {
			h.app.TelegramBot.SendMessage(userDB.ID, text, &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			})
		}
	}
}
