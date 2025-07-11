package handlers

import (
	"bot/internal/common/database"

	"github.com/bwmarrin/discordgo"
)

func (h *handlers) GuildAddHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	var existingGuild database.Guild
	if err := h.db.First(&existingGuild, &database.Guild{ID: g.ID}).Error; err == nil {
		return
	}

	guild := database.Guild{ID: g.ID, Name: g.Name}

	for _, channel := range g.Guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			guild.Channels = append(guild.Channels, &database.Channel{ID: channel.ID, Name: channel.Name})
		}
	}

	h.db.Create(&guild)
}

func (h *handlers) GuildUpdateHandler(s *discordgo.Session, g *discordgo.GuildUpdate) {
	h.db.Model(&database.Guild{}).Where(&database.Guild{ID: g.ID}).Update("name", g.Name)
}

func (h *handlers) GuildDeleteHandler(s *discordgo.Session, g *discordgo.GuildDelete) {
	h.db.Unscoped().Where(&database.Guild{ID: g.ID}).Delete(&database.Guild{})
}
