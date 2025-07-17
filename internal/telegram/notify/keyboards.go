package notify

import (
	"fmt"

	"bot/internal/common/database"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func guildsMarkup(guilds []database.Guild) gotgbot.InlineKeyboardMarkup {
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

func channelsMarkup(channels []database.Channel, userChannels []database.Channel) gotgbot.InlineKeyboardMarkup {
	var buttons [][]gotgbot.InlineKeyboardButton
	var row []gotgbot.InlineKeyboardButton

	for i, c := range channels {
		text := c.Name
		inList := false
		for _, uc := range userChannels {
			if uc.ID == c.ID {
				inList = true
				break
			}
		}
		if inList {
			text += " [🗸]"
		} else {
			text += " [✘]"
		}

		button := gotgbot.InlineKeyboardButton{
			Text:         text,
			CallbackData: fmt.Sprintf("channel:%s:%s", c.GuildID, c.ID),
		}

		row = append(row, button)

		if len(row) == 2 {
			buttons = append(buttons, row)
			row = []gotgbot.InlineKeyboardButton{}
		}

		if i == len(channels)-1 && len(row) > 0 {
			buttons = append(buttons, row)
		}
	}

	buttons = append(buttons, []gotgbot.InlineKeyboardButton{{Text: "◀️", CallbackData: "channel:back:0"}})

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}
