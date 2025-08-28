package handlers

import (
	"bot/internal/discord/app"

	"github.com/bwmarrin/discordgo"
)

type handlers struct {
	app      *app.AppContext
	commands []*discordgo.ApplicationCommand
}

func New(app *app.AppContext, commands []*discordgo.ApplicationCommand) *handlers {
	return &handlers{app, commands}
}

func (h *handlers) Handlers() []interface{} {
	hdls := []interface{}{
		h.GuildAddHandler,
		h.GuildUpdateHandler,
		h.GuildDeleteHandler,

		h.VoiceJoinHandler,

		h.ChannelAddHandler,
		h.ChannelUpdateHandler,
		h.ChannelDeleteHandler,
	}
	return hdls
}
