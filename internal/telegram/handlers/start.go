package handlers

import (
	"strings"

	"bot/internal/common/database"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

func (h *handlers) StartHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) > 1 {
		user, _ := ctx.Data["user"].(*database.User)

		data := strings.Split(args[1], "_")

		var dbGuild database.Guild
		if err := h.db.Preload("Channels").First(&dbGuild, "id = ?", data[0]).Error; err != nil {
			h.logger.Error("failed to get guild", zap.Error(err))
			return nil
		}

		h.db.Model(&user).Association("Guilds").Append(&dbGuild)
		h.db.Model(&user).Association("Channels").Append(&dbGuild.Channels)
		h.db.Save(&user)

		b.SendMessage(ctx.EffectiveChat.Id, "Success added guild: "+dbGuild.Name, nil)
	} else {
		b.SendMessage(ctx.EffectiveChat.Id, "Hello "+ctx.EffectiveUser.FirstName, nil)
	}
	return nil
}
