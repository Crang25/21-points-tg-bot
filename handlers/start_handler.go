package handlers

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Crang25/21points_tgbot/keyboards"
	"github.com/Crang25/21points_tgbot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ❤️ ♦️ ♣ ♠
// Deck - колода карт
var deck = [36]models.Card{
	{Worth: 6, Value: "|❤️ 6|"}, {Worth: 7, Value: "|❤️ 7|"}, {Worth: 8, Value: "|❤️ 8|"},
	{Worth: 9, Value: "|❤️ 9|"}, {Worth: 10, Value: "|❤️ 10|"}, {Worth: 2, Value: "|❤️ В|"},
	{Worth: 3, Value: "|❤️ Д|"}, {Worth: 4, Value: "|❤️ К|"}, {Worth: 11, Value: "|❤️ Т|"},

	{Worth: 6, Value: "|♦️ 6|"}, {Worth: 7, Value: "|♦️ 7|"}, {Worth: 8, Value: "|♦️ 8|"},
	{Worth: 9, Value: "|♦️ 9|"}, {Worth: 10, Value: "|♦️ 10|"}, {Worth: 2, Value: "|♦️ В|"},
	{Worth: 3, Value: "|♦️ Д|"}, {Worth: 4, Value: "|♦️ К|"}, {Worth: 11, Value: "|♦️ Т|"},

	{Worth: 6, Value: "|♣ 6|"}, {Worth: 7, Value: "|♣ 7|"}, {Worth: 8, Value: "|♣ 8|"},
	{Worth: 9, Value: "|♣ 9|"}, {Worth: 10, Value: "|♣ 10|"}, {Worth: 2, Value: "|♣ В|"},
	{Worth: 3, Value: "|♣ Д|"}, {Worth: 4, Value: "|♣ К|"}, {Worth: 11, Value: "|♣ Т|"},

	{Worth: 6, Value: "|♠ 6|"}, {Worth: 7, Value: "|♠ 7|"}, {Worth: 8, Value: "|♠ 8|"},
	{Worth: 9, Value: "|♠ 9|"}, {Worth: 10, Value: "|♠ 10|"}, {Worth: 2, Value: "|♠ В|"},
	{Worth: 3, Value: "|♠ Д|"}, {Worth: 4, Value: "|♠ К|"}, {Worth: 11, Value: "|♠ Т|"},
}

// Счетчик
var counter = len(deck) - 1

// Start - игра
func Start(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, chatID int64, firstName string) {
	// Перемешываем колоду
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	//  Карты игрока и дилера соответственно
	player := models.PlayerDeck{}
	banker := models.PlayerDeck{}

	// Первая карта дилера
	banker.PointsQuan += deck[counter].Worth
	banker.Deck = append(banker.Deck, deck[counter].Value)
	counter--

	// Первая карта игрока
	player.PointsQuan = deck[counter].Worth
	player.Deck = append(player.Deck, deck[counter].Value)
	counter--

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("Игра началась!🎰\nКол-во очков дилера: %d 💵\nКарты дилера: %s\n",
			banker.PointsQuan, strCards(banker.Deck)))

	bot.Send(msg)

	msg.Text = fmt.Sprintf("%s, кол-во ваших очков %d 💵\nВаши карты: %s",
		firstName, player.PointsQuan, strCards(player.Deck))
	msg.ReplyMarkup = keyboards.GameKeyboard
	counter--

	bot.Send(msg)

	// Принимаем сообщения от пользователя
	for update := range updates {

		if update.Message == nil {
			continue
		}

		// Сдача карт игроку, пока он сам не остановит, либо пока не наберется больше либо равно 21
		if update.Message.Text == "Завершить игру ❌" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Игра завершена")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg.ReplyMarkup = keyboards.AgainKeyboard
			bot.Send(msg)
			return
		} else if update.Message.Text == "Еще 👍" { // Даем карту игроку

			player.PointsQuan += deck[counter].Worth
			player.Deck = append(player.Deck, deck[counter].Value)
			counter--

			msg.Text = fmt.Sprintf("%s, кол-во ваших очков %d 💵\nВаши карты: %s",
				firstName, player.PointsQuan, strCards(player.Deck))
			msg.ReplyMarkup = keyboards.GameKeyboard
			bot.Send(msg)

		}

		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		// Игрок набрал 21 очко
		if player.PointsQuan == 21 {
			msg.Text = fmt.Sprintf("%s, вы выйграли!🏆\nУ вас 21 очко\nВаши карты: %s",
				firstName, strCards(player.Deck))
			endGame(&msg, bot)
			return

			// Игрок проиграл - набрал больше 21
		} else if player.PointsQuan >= 21 {
			msg.Text = fmt.Sprintf("%s, вы проиграли.\nУ вас %d очка(ов), это больше 21\nВаши карты: %s",
				firstName, player.PointsQuan, strCards(player.Deck))
			endGame(&msg, bot)
			return

		} else if update.Message.Text == "Хватит 👎" {
			generateBankerDeck(&banker)
			// Если выйграл банкир
			if banker.IsWin || banker.PointsQuan > player.PointsQuan {
				msg.Text = fmt.Sprintf(
					"%s, вы проиграли\nКол-во ваших очков: %d 💵\nВаши карты: %s\n\nКол-во очков дилера: %d 💵\nКарты дилера: %s",
					firstName, player.PointsQuan, player.Deck, banker.PointsQuan, banker.Deck)
				endGame(&msg, bot)
				return
				// Выйграл игрок
			} else if banker.PointsQuan < player.PointsQuan {
				msg.Text = fmt.Sprintf(
					"%s, вы выйграли! 💰\nКол-во ваших очков: %d 💵\nВаши карты: %s\n\nКол-во очков дилера: %d 💵\nКарты дилера: %s",
					firstName, player.PointsQuan, player.Deck, banker.PointsQuan, banker.Deck)
				endGame(&msg, bot)
				return
			}
		}
	}

}

// Функция, "склеивающая" слайс строк в одну строку
func strCards(cards []string) string {
	str := ""
	for _, card := range cards {
		str += card + ", "
	}
	return str
}

// Функция, в которой карты раздаются дилеру, пока сумма очков меньше 17
func generateBankerDeck(bankir *models.PlayerDeck) {
	for bankir.PointsQuan < 21 {
		if bankir.PointsQuan == 21 {
			bankir.IsWin = true
			return
		} else if bankir.PointsQuan < 21 && bankir.PointsQuan >= 17 {
			return
		}

		bankir.PointsQuan += deck[counter].Worth
		bankir.Deck = append(bankir.Deck, deck[counter].Value)
		counter--
	}
}

// Завершние раунда
func endGame(msg *tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyMarkup = keyboards.AgainKeyboard
	bot.Send(msg)
}
