package handlers

import (
	"bot/internal/database"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *handlers) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) > 1 {
		user, _ := ctx.Data["user"].(*database.User)

		data := strings.Split(args[1], "_")

		var dbGuild database.Guild
		h.db.Preload("Channels").First(&dbGuild, data[0])

		h.db.Model(&user).Association("Guilds").Append(&dbGuild)
		h.db.Model(&user).Association("Channels").Append(&dbGuild.Channels)

		b.SendMessage(ctx.EffectiveChat.Id, "Success added guild: "+dbGuild.Name, nil)
	} else {
		b.SendMessage(ctx.EffectiveChat.Id, "Hello "+ctx.EffectiveUser.FirstName, nil)
	}
	return nil
}
