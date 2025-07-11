package keyboards

import (
	"fmt"

	"bot/internal/common/database"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func GuildsMarkup(guilds []*database.Guild) gotgbot.InlineKeyboardMarkup {
	var buttons [][]gotgbot.InlineKeyboardButton
	var row []gotgbot.InlineKeyboardButton

	for _, g := range guilds {
		row = append(row, gotgbot.InlineKeyboardButton{Text: g.Name, CallbackData: fmt.Sprintf("guild:%s", g.ID)})
		if len(row) == 2 {
			buttons = append(buttons, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		buttons = append(buttons, row)
	}

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}
