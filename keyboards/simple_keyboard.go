package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// GameKeyboard - Клавиатура, отображающаяся во время игры
var GameKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Еще 👍"),
		tgbotapi.NewKeyboardButton("Хватит 👎"),
		tgbotapi.NewKeyboardButton("Завершить игру ❌"),
	),
)
