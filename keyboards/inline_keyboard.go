package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// AgainKeyboard - инлайн клавиатура, предлагающая начать игру занова
var AgainKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Начать заново", "again"),
	),
)

// RulesKeyboard - инлайн клавиатура, с предложениями посмотреть правила игры
var RulesKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("Статья",
			"https://homeofgames.ru/game/21-ochko"),
		tgbotapi.NewInlineKeyboardButtonURL("Видео",
			"https://www.youtube.com/watch?v=YxYx0eVeBK0"),
	),
)
