package handlers

import (
	"strings"

	"bot/internal/common/database"
	"bot/internal/telegram/keyboards"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *handlers) NotifyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	b.SendMessage(ctx.EffectiveChat.Id, "Notify:", &gotgbot.SendMessageOpts{ReplyMarkup: keyboards.GuildsMarkup(user.Guilds)})

	return nil
}

func (h *handlers) NotifyGuildHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	data := strings.Split(ctx.CallbackQuery.Data, ":")

	var guild *database.Guild
	h.db.Preload("Channels").First(&guild, "id = ?", data[1])

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: keyboards.ChannelsMarkup(guild.Channels, user.Channels),
	})

	return nil
}

func (h *handlers) NotifyChannelHandler(b *gotgbot.Bot, ctx *ext.Context) error {
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
	h.db.Preload("Channels").First(&guild, "id = ?", guildId)
	var channel database.Channel
	h.db.First(&channel, "id = ?", channelId)

	var inUser bool
	for _, c := range user.Channels {
		if c.ID == channel.ID {
			inUser = true
			h.db.Model(&user).Association("Channels").Delete(&channel)
			break
		}
	}
	if !inUser {
		h.db.Model(&user).Association("Channels").Append(&channel)
	}

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: keyboards.ChannelsMarkup(guild.Channels, user.Channels),
	})

	return nil
}
