package middlewares

import (
	"bot/internal/common/database"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (m *middlewares) UserMiddleware(_ *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	if user == nil {
		return nil
	}

	var dbUser database.User
	m.db.Preload("Guilds").Preload("Channels").FirstOrCreate(&dbUser, database.User{TelegramID: user.Id})

	ctx.Data["user"] = &dbUser

	return nil
}
