package commander

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	"bot/config"
)

type CmdHandler func(string) string

type Commander struct {
	bot   *tgbotapi.BotAPI
	route map[string]CmdHandler
}

func Init() (*Commander, error) {
	log.Println("init commander")
	bot, err := tgbotapi.NewBotAPI(config.GetApiKey())
	if err != nil {
		return nil, err
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Commander{
		bot:   bot,
		route: make(map[string]CmdHandler),
	}, errors.Wrap(err, "er")
}

func (c *Commander) RegisterCommander(cmd string, f CmdHandler) {
	c.route[cmd] = f
}

func (c *Commander) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			c.handleMessage(update.Message)
		}
	}
}

func (c *Commander) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	if cmd := message.Command(); cmd != "" {
		if f, ok := c.route[cmd]; ok {
			msg.Text = f(message.CommandArguments())
		} else {
			msg.Text = fmt.Sprintf("command not found <%s>", cmd)
		}
	} else {
		msg.Text = fmt.Sprintf("your message <%s>", message.Text)
	}
	_, err := c.bot.Send(msg)
	if err != nil {
		log.Printf("answer error: %v\n", err)
	}
}
