package commander

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Commander struct {
	bot *tgbotapi.BotAPI
}

func Init(bot *tgbotapi.BotAPI) (*Commander, error) {
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Commander{
		bot: bot,
	}, nil
}

func (c *Commander) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			res := fmt.Sprintf("you send <%v>", update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
			msg.ReplyToMessageID = update.Message.MessageID
			c.bot.Send(msg)
		}
	}
}
