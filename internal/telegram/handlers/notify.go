package handlers

import (
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/telegram/keyboards"
)

func (h *Handlers) NotifyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	b.SendMessage(ctx.EffectiveChat.Id, "Notify:", &gotgbot.SendMessageOpts{ReplyMarkup: keyboards.GuildsMarkup(user.Guilds)})

	return nil
}

func (h *Handlers) NotifyGuildHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	data := strings.Split(ctx.CallbackQuery.Data, ":")

	var guild *database.Guild
	h.DB.Preload("Channels").First(&guild, "id = ?", data[1])

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: keyboards.ChannelsMarkup(guild.Channels, user.Channels),
	})

	return nil
}

func (h *Handlers) NotifyChannelHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	data := strings.Split(ctx.CallbackQuery.Data, ":")
	guildId, channelId := data[1], data[2]

	if guildId == "back" {
		b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
			ChatId:      ctx.EffectiveChat.Id,
			MessageId:   ctx.EffectiveMessage.MessageId,
			ReplyMarkup: keyboards.GuildsMarkup(user.Guilds),
		})
		return nil
	}

	var guild *database.Guild
	h.DB.Preload("Channels").First(&guild, "id = ?", guildId)
	var channel database.Channel
	h.DB.First(&channel, "id = ?", channelId)

	var inUser bool
	for _, c := range user.Channels {
		if c.ID == channel.ID {
			inUser = true
			h.DB.Model(&user).Association("Channels").Delete(&channel)
			break
		}
	}
	if !inUser {
		h.DB.Model(&user).Association("Channels").Append(&channel)
	}

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: keyboards.ChannelsMarkup(guild.Channels, user.Channels),
	})

	return nil
}
