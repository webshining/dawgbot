package middlewares

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/webshining/internal/common/database"
)

func (m *Middlewares) UserMiddleware(_ *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	if user == nil {
		return nil
	}

	var dbUser database.User
	m.DB.Preload("Guilds").Preload("Channels").FirstOrCreate(&dbUser, database.User{ID: user.Id})

	ctx.Data["user"] = &dbUser

	return nil
}
