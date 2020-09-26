package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Crang25/21points_tgbot/logger"

	"github.com/Crang25/21points_tgbot/handlers"
	"github.com/Crang25/21points_tgbot/keyboards"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/joho/godotenv"
)

// TOKEN - токен бота
var TOKEN string

func init() {
	// Загружаем переменные окружения
	err := godotenv.Load()
	logger.CheckErr(err)

	// Загружаем токен бота
	var isExists bool
	TOKEN, isExists = os.LookupEnv("TOKEN")
	if !isExists {
		logger.CheckErr(errors.New("failed to load TOKEN"))
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		logger.CheckErr(err)
	}

	// Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Объявляем канал, куда будут приходить обновления
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.CheckErr(err)
	}

	// Принимаем обновления/сообщения
	for update := range updates {

		if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "again" {
				handlers.Start(updates, bot, update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.From.FirstName)
			}
		}

		if update.Message != nil {

			// Обработка команд, присланных боту
			if update.Message.IsCommand() {

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

				switch update.Message.Command() {
				case "start":
					handlers.Start(updates, bot, update.Message.Chat.ID, update.Message.From.FirstName)
				case "stop":
					msg.Text = "Игра завершена"
					msg.ReplyMarkup = keyboards.AgainKeyboard
				case "help":
					// TODO: сделать handler для help
					msg.Text = `Бот, с которым можно поиграть в игру "21 очко" или просто "Очко".
	Чтобы начать новую игру отправьте боту /start.
	Чтобы увидеть список всех команд бота, отправьте ему /commands.
	С общими правилами игры можете ознакомиться, отправив боту команду /rules`
				case "commands":
					msg.Text =
						fmt.Sprintf(
							"Список команд:\n/start - начать новую игру\n/stop - закончить текущую игру\n/help - посмотреть базовую информацию о работе с ботом\n/commands - посмотреть все команды бота\n/rules - ознакомиться с правилами игры\n/close_keyboard - убрать клавиатуру")
				case "rules":
					msg.Text = "Правила игры"
					msg.ReplyMarkup = keyboards.RulesKeyboard
				case "close_keyboard":
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				default:
					msg.Text = "Такой команды нет. Отправьте /commands, чтобы увидеть список всех команд"
				}

				bot.Send(msg)
			}

			if update.Message.Text == "Завершить игру ❌" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Игра завершена")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				msg.ReplyMarkup = keyboards.AgainKeyboard
				bot.Send(msg)
			}
		}

	}

}
