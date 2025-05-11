package handlers

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/telegram/keyboards"
)

func (h *Handlers) NotifyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, ok := ctx.Data["user"].(*database.User)
	if !ok {
		fmt.Println("error get user")
		return nil
	}

	b.SendMessage(ctx.EffectiveChat.Id, "Notify:", &gotgbot.SendMessageOpts{ReplyMarkup: keyboards.GuildsMarkup(user.Guilds)})

	return nil
}

func (h *Handlers) NotifyGuildHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	return nil
}

func (h *Handlers) NotifyChannelHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	return nil
}
