package handlers

import (
	"bot/internal/common/database"

	"github.com/bwmarrin/discordgo"
)

func (h *handlers) ChannelAddHandler(s *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Type != discordgo.ChannelTypeGuildVoice {
		return
	}

	var existingChannel database.Channel
	if err := h.db.First(&existingChannel, &database.Channel{ID: c.ID, GuildID: c.GuildID}).Error; err == nil {
		return
	}

	h.db.Create(&database.Channel{ID: c.ID, Name: c.Name, GuildID: c.GuildID})
}

func (h *handlers) ChannelUpdateHandler(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	if c.Type != discordgo.ChannelTypeGuildVoice {
		return
	}

	h.db.Model(&database.Channel{}).Where(&database.Channel{ID: c.ID, GuildID: c.GuildID}).Update("name", c.Name)
}

func (h *handlers) ChannelDeleteHandler(s *discordgo.Session, c *discordgo.ChannelDelete) {
	if c.Type != discordgo.ChannelTypeGuildVoice {
		return
	}

	h.db.Unscoped().Where(&database.Channel{ID: c.ID, GuildID: c.GuildID}).Delete(&database.Channel{})
}
