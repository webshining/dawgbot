package notify

import (
	"bot/internal/common/database"
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

func (n *notify) notifyHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	var guilds []database.Guild
	n.app.DB.Model(&user).Association("Guilds").Find(&guilds)
	ctx.EffectiveMessage.Reply(b, "Notify:", &gotgbot.SendMessageOpts{ReplyMarkup: guildsMarkup(guilds)})

	return nil
}

func (n *notify) notifyGuildHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	data := strings.Split(ctx.CallbackQuery.Data, ":")

	var guildChannels []database.Channel
	n.app.DB.Model(&database.Guild{ID: data[1]}).Association("Channels").Find(&guildChannels)

	var userChannels []database.Channel
	n.app.DB.Model(&user).Association("Channels").Find(&userChannels)

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: channelsMarkup(guildChannels, userChannels),
	})

	return nil
}

func (n *notify) notifyChannelHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user, _ := ctx.Data["user"].(*database.User)

	data := strings.Split(ctx.CallbackQuery.Data, ":")
	guildId, channelId := data[1], data[2]

	if guildId == "back" {
		var guilds []database.Guild
		n.app.DB.Model(&user).Association("Guilds").Find(&guilds)
		ctx.EffectiveMessage.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{ReplyMarkup: guildsMarkup(guilds)})

		return nil
	}

	var guildChannels []database.Channel
	n.app.DB.Model(&database.Guild{ID: guildId}).Association("Channels").Find(&guildChannels)
	var userChannels []database.Channel
	n.app.DB.Model(&user).Association("Channels").Find(&userChannels)
	var channel database.Channel
	n.app.DB.First(&channel, channelId)

	var inUser bool
	for i, c := range userChannels {
		if c.ID == channel.ID {
			inUser = true
			n.app.DB.Model(&user).Association("Channels").Delete(&channel)

			userChannels = append(userChannels[:i], userChannels[i+1:]...)
			break
		}
	}
	if !inUser {
		n.app.DB.Model(&user).Association("Channels").Append(&channel)
		userChannels = append(userChannels, channel)
	}

	fmt.Printf("%#v", guildChannels)

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: channelsMarkup(guildChannels, userChannels),
	})

	return nil
}

func (n *notify) Handlers(dp *ext.Dispatcher, group int) {
	dp.AddHandlerToGroup(handlers.NewCommand("notify", n.notifyHandler), group)
	dp.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("guild"), n.notifyGuildHandler), group)
	dp.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("channel"), n.notifyChannelHandler), group)
}
