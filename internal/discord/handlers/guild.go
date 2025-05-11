package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/webshining/internal/common/database"
)

func (h *Handlers) GuildAddHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	var existingGuild database.Guild
	if err := h.DB.First(&existingGuild, &database.Guild{ID: g.ID}).Error; err == nil {
		return
	}

	guild := database.Guild{ID: g.ID, Name: g.Name}

	for _, channel := range g.Guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			guild.Channels = append(guild.Channels, &database.Channel{ID: channel.ID, Name: channel.Name})
		}
	}

	h.DB.Create(&guild)
}

func (h *Handlers) GuildUpdateHandler(s *discordgo.Session, g *discordgo.GuildUpdate) {
	h.DB.Model(&database.Guild{}).Where(&database.Guild{ID: g.ID}).Update("name", g.Name)
}

func (h *Handlers) GuildDeleteHandler(s *discordgo.Session, g *discordgo.GuildDelete) {
	h.DB.Unscoped().Where(&database.Guild{ID: g.ID}).Delete(&database.Guild{})
}
