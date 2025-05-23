package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/webshining/internal/common/database"
	"go.uber.org/zap"
)

func (h *Handlers) StartHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) > 1 {
		user, _ := ctx.Data["user"].(*database.User)

		var dbGuild database.Guild
		if err := h.DB.Preload("Channels").First(&dbGuild, "id = ?", args[1]).Error; err != nil {
			h.logger.Error("failed to get guild", zap.Error(err))
			return nil
		}

		h.DB.Model(&user).Association("Guilds").Append(&dbGuild)
		h.DB.Model(&user).Association("Channels").Append(&dbGuild.Channels)

		b.SendMessage(ctx.EffectiveChat.Id, "Success added guild: "+dbGuild.Name, nil)
	} else {
		b.SendMessage(ctx.EffectiveChat.Id, "Hello "+ctx.EffectiveUser.FirstName, nil)
	}
	return nil
}
