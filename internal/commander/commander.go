package commander

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kshmatov/masterclass1/config"
)

var UnknownCommand = errors.New("unknown command")

type Commander struct {
	bot *tgbotapi.BotAPI
}

func Init() (*Commander, error) {
	bot, err := tgbotapi.NewBotAPI(config.ApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "init tgbot")
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Commander{
		bot: bot,
	}, nil
}

func (c *Commander) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			res := fmt.Sprintf("you send <%v>", update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
			msg.ReplyToMessageID = update.Message.MessageID
			_, err := c.bot.Send(msg)
			if err != nil {
				return errors.Wrap(err, "send tg message")
			}
		}
	}
	return nil
}
