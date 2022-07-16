package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/kshmatov/masterclass1/config"
	"gitlab.ozon.dev/kshmatov/masterclass1/internal/commander"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.ApiKey)
	if err != nil {
		log.Panic(err)
	}
	cmd, err := commander.Init(bot)
	if err != nil {
		log.Panic(err)
	}
	cmd.Run()
}
