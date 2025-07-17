package user

import (
	"bot/internal/common/database"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func (u *user) userMiddleware(_ *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	if user == nil {
		return nil
	}

	var dbUser database.User
	u.app.DB.FirstOrCreate(&dbUser, database.User{ID: user.Id})

	ctx.Data["user"] = &dbUser

	return nil
}

func (u *user) Middlewares(dp *ext.Dispatcher, group int) {
	dp.AddHandlerToGroup(handlers.NewMessage(message.Command, u.userMiddleware), group)
	dp.AddHandlerToGroup(handlers.NewCallback(callbackquery.All, u.userMiddleware), group)
}
