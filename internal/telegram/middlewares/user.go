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
	if err := m.DB.Preload("Guilds").Preload("Channels").Preload("Guilds.Channels").First(&dbUser, "id = ?", user.Id).Error; err != nil {
		m.DB.Create(&database.User{ID: user.Id})
	}

	ctx.Data["user"] = &dbUser

	return nil
}
