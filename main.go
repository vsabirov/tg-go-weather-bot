package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vsabirov/tggoweatherbot/cmdhandlers"
	"github.com/vsabirov/tggoweatherbot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getCriticalEnvValue(name string, description string) string {
	value, valueExists := os.LookupEnv(name)
	if !valueExists {
		log.Panicln(description, "not found. Please specify it through the", name, "environment variable.")
	}

	return value
}

func main() {
	// Build the database credentials from environment variables.
	port, err := strconv.Atoi(getCriticalEnvValue("DB_PORT", "Database port"))
	if err != nil {
		log.Panicln("Environment variable DB_PORT should be a number.")
	}

	dbCredentials := database.Credentials{}
	dbCredentials.Host = getCriticalEnvValue("DB_HOST", "Database hostname")
	dbCredentials.Port = port
	dbCredentials.User = getCriticalEnvValue("DB_USER", "Database username")
	dbCredentials.Password = getCriticalEnvValue("DB_PASSWORD", "Database password")
	dbCredentials.Database = getCriticalEnvValue("DB_NAME", "Database name")

	err = database.Connect(dbCredentials)
	if err != nil {
		log.Panicln("Failed to connect to the database: ", err)
	}

	log.Println("Connected to the database.")

	defer database.Disconnect()

	// Build the bot from environment variables.
	bot, err := tgbotapi.NewBotAPI(getCriticalEnvValue("BOT_TOKEN", "Bot token"))
	if err != nil {
		log.Panicln(err)
	}

	bot.Debug = os.Getenv("PRODUCTION") != "true"

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	// List of functions that will react and respond to certain user messages. All messages should be lowercase.
	handlers := make(map[string]cmdhandlers.CommandHandler)

	handlers["начать"] = cmdhandlers.Start
	handlers["информация"] = cmdhandlers.Info
	handlers["статистика"] = cmdhandlers.Stats

	// Bot event loop.
	for update := range updates {
		if update.Message != nil {
			payload := update.Message.Text
			arguments := strings.Split(payload, " ")

			command := strings.ToLower(arguments[0])

			if handlers[command] != nil {
				go handlers[command](update, bot)
			} else {
				go handlers["начать"](update, bot)
			}
		}
	}
}
