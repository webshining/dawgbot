package telegram

import "github.com/PaulSonOfLars/gotgbot/v2"

func (b *Bot) SendMessage(chatId int64, text string, opts *gotgbot.SendMessageOpts) {
	b.bot.SendMessage(chatId, text, opts)
}

func (b *Bot) SendPhoto(chatId int64, photo gotgbot.InputFileOrString, opts *gotgbot.SendPhotoOpts) {
	b.bot.SendPhoto(chatId, photo, opts)
}
