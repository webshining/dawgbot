package handlers

import (
	"bot/internal/database"
	"bot/internal/telegram/keyboards"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *handlers) Notify(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	var guilds []database.Guild
	h.db.Model(&user).Association("Guilds").Find(&guilds)
	ctx.EffectiveMessage.Reply(b, "Notify:", &gotgbot.SendMessageOpts{ReplyMarkup: keyboards.GuildsMarkup(guilds)})

	return nil
}

func (h *handlers) NotifyGuild(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	data := strings.Split(ctx.CallbackQuery.Data, ":")

	var guildChannels []database.Channel
	h.db.Model(&database.Guild{ID: data[1]}).Association("Channels").Find(&guildChannels)

	var userChannels []database.Channel
	h.db.Model(&user).Association("Channels").Find(&userChannels)

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: keyboards.ChannelsMarkup(guildChannels, userChannels),
	})

	return nil
}

func (h *handlers) NotifyChannel(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	data := strings.Split(ctx.CallbackQuery.Data, ":")
	guildId, channelId := data[1], data[2]

	if guildId == "back" {
		var guilds []database.Guild
		h.db.Model(&user).Association("Guilds").Find(&guilds)
		ctx.EffectiveMessage.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{ReplyMarkup: keyboards.GuildsMarkup(guilds)})

		return nil
	}

	var guildChannels []database.Channel
	h.db.Model(&database.Guild{ID: guildId}).Association("Channels").Find(&guildChannels)
	var userChannels []database.Channel
	h.db.Model(&user).Association("Channels").Find(&userChannels)
	var channel database.Channel
	h.db.First(&channel, channelId)

	var inUser bool
	for i, c := range userChannels {
		if c.ID == channel.ID {
			inUser = true
			h.db.Model(&user).Association("Channels").Delete(&channel)

			userChannels = append(userChannels[:i], userChannels[i+1:]...)
			break
		}
	}
	if !inUser {
		h.db.Model(&user).Association("Channels").Append(&channel)
		userChannels = append(userChannels, channel)
	}

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: keyboards.ChannelsMarkup(guildChannels, userChannels),
	})

	return nil
}
