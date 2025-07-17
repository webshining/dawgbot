package start

import (
	"bot/internal/common/database"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func (s *start) startHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	args := ctx.Args()
	if len(args) > 1 {
		user, _ := ctx.Data["user"].(*database.User)

		data := strings.Split(args[1], "_")

		var dbGuild database.Guild
		s.app.DB.Preload("Channels").First(&dbGuild, data[0])

		s.app.DB.Model(&user).Association("Guilds").Append(&dbGuild)
		s.app.DB.Model(&user).Association("Channels").Append(&dbGuild.Channels)

		b.SendMessage(ctx.EffectiveChat.Id, "Success added guild: "+dbGuild.Name, nil)
	} else {
		b.SendMessage(ctx.EffectiveChat.Id, "Hello "+ctx.EffectiveUser.FirstName, nil)
	}
	return nil
}

func (s *start) Handlers(dp *ext.Dispatcher, group int) {
	dp.AddHandlerToGroup(handlers.NewCommand("start", s.startHandler), group)
}
