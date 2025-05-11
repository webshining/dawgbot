package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/webshining/internal/common/database"
)

func (h *Handlers) StartHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) > 1 {
		user, ok := ctx.Data["user"].(*database.User)
		if !ok {
			fmt.Println("error get user")
			return nil
		}

		var dbGuild database.Guild
		if err := h.DB.Preload("Channels").First(&dbGuild, "id = ?", args[1]).Error; err != nil {
			fmt.Println("error get guild")
			return nil
		}

		user.Guilds = append(user.Guilds, &dbGuild)
		user.Channels = dbGuild.Channels
		h.DB.Save(user)

		b.SendMessage(ctx.EffectiveChat.Id, "Success added guild: "+dbGuild.Name, nil)
	} else {
		b.SendMessage(ctx.EffectiveChat.Id, "Hello "+ctx.EffectiveUser.FirstName, nil)
	}
	return nil
}
