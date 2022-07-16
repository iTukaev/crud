package commander

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kshmatov/masterclass1/config"
	"gitlab.ozon.dev/kshmatov/masterclass1/internal/storage"
)

const (
	listCmd = "list"
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
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if cmd := update.Message.Command(); cmd != "" {
			switch cmd {
			case listCmd:
				msg.Text = listFunc()
			default:
				msg.Text = UnknownCommand.Error()
			}
		} else {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg.Text = fmt.Sprintf("you send <%v>", update.Message.Text)
		}
		_, err := c.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "send tg message")
		}
	}
	return nil
}

func listFunc() string {
	data := storage.List()
	res := make([]string, 0, len(data))
	for _, v := range data {
		res = append(res, v.String())
	}
	return strings.Join(res, "\n")
}
