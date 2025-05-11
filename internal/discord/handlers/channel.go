package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/webshining/internal/common/database"
)

func (h *Handlers) ChannelAddHandler(s *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Type != discordgo.ChannelTypeGuildVoice {
		return
	}

	var existingChannel database.Channel
	if err := h.DB.First(&existingChannel, &database.Channel{ID: c.ID, GuildID: c.GuildID}).Error; err == nil {
		return
	}

	h.DB.Create(&database.Channel{ID: c.ID, Name: c.Name, GuildID: c.GuildID})
}

func (h *Handlers) ChannelUpdateHandler(s *discordgo.Session, c *discordgo.ChannelUpdate) {
	if c.Type != discordgo.ChannelTypeGuildVoice {
		return
	}

	h.DB.Model(&database.Channel{}).Where(&database.Channel{ID: c.ID, GuildID: c.GuildID}).Update("name", c.Name)
}

func (h *Handlers) ChannelDeleteHandler(s *discordgo.Session, c *discordgo.ChannelDelete) {
	if c.Type != discordgo.ChannelTypeGuildVoice {
		return
	}

	h.DB.Unscoped().Where(&database.Channel{ID: c.ID, GuildID: c.GuildID}).Delete(&database.Channel{})
}
