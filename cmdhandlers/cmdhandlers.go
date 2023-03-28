package cmdhandlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/vsabirov/tggoweatherbot/executors"
)

// User request handler function type.
type CommandHandler func(update tgbotapi.Update, bot *tgbotapi.BotAPI)

// Welcome message for new users.
func Start(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	userKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Информация Москва"),
			tgbotapi.NewKeyboardButton("Статистика"),
		),
	)

	response := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Введите 'Информация <Название города>', или 'Статистика'. Например 'Информация Москва'")

	response.ReplyMarkup = userKeyboard

	bot.Send(response)
}

// This function handles an information request from the user.
func Info(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	response.ReplyToMessageID = update.Message.MessageID

	arguments := strings.Split(update.Message.Text, " ")
	if len(arguments) < 2 {
		response.Text = "Пожалуйста, введите название города после слова 'Информация' через пробел."

		bot.Send(response)
		return
	}

	city := arguments[1]

	geodata, err := executors.FetchGeoInformation(city)
	if err != nil {
		log.Println(err)

		response.Text = "Произошла ошибка при поиске координат города."

		bot.Send(response)
		return
	}

	weather, err := executors.FetchWeatherInformation(geodata)
	if err != nil {
		log.Println(err)

		response.Text = "Произошла ошибка при поиске информации о погоде."

		bot.Send(response)
		return
	}

	response.Text = fmt.Sprintf(
		"🏙️ Погода в %s, %s:\n🌡️ Температура: %.2f°C\n💧 Дождливость: %.1f%%\n‍💨 Скорость ветра: %.2f км/ч",
		geodata.City, geodata.Country,
		weather.Temperature, weather.RainLevel*100, weather.WindSpeed)

	bot.Send(response)

	go func() {
		err = executors.RecordStats(update.Message.From.ID, geodata.City)
		if err != nil {
			log.Println(err)
		}
	}()
}

// This function handles a statistics request from the user.
func Stats(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	response := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	response.ReplyToMessageID = update.Message.MessageID

	stats, err := executors.FetchUserStats(update.Message.From.ID)
	if err != nil {
		log.Println(err)

		response.Text = "Произошла ошибка при попытке получить вашу статистику. Возможно, вы ещё не делали запросов у нас."

		bot.Send(response)
		return
	}

	response.Text = fmt.Sprintf("🕒 Первый запрос: %s\n🧮 Количество запросов: %d\n🏙️ Последний город: %s",
		stats.FirstRequest.UTC().Format(time.UnixDate), stats.AmountOfRequests,
		stats.LastCity)

	bot.Send(response)
}
