package keyboards

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/webshining/internal/common/database"
)

func GuildsMarkup(guilds []*database.Guild) gotgbot.InlineKeyboardMarkup {
	var buttons [][]gotgbot.InlineKeyboardButton
	for _, g := range guilds {
		row := []gotgbot.InlineKeyboardButton{
			{Text: g.Name, CallbackData: fmt.Sprintf("guild:%s", g.ID)},
		}
		buttons = append(buttons, row)
	}
	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}
