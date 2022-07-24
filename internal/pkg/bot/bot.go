package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
)

type Interface interface {
	RegisterCommander(cmd commandPkg.Interface)
	Run()
}

func MustNew(id string) Interface {
	log.Println("Init bot")
	bot, err := tgbotapi.NewBotAPI(id)
	if err != nil {
		log.Panic(errors.Wrap(err, "new API bot"))
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &commander{
		bot:   bot,
		route: make(map[string]commandPkg.Interface),
	}
}

type commander struct {
	bot   *tgbotapi.BotAPI
	route map[string]commandPkg.Interface
}

// RegisterCommander - not thread safe
func (c *commander) RegisterCommander(cmd commandPkg.Interface) {
	c.route[cmd.Name()] = cmd
}

func (c *commander) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			c.handleMessage(update.Message)
		}
	}
}

func (c *commander) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	if cmdName := message.Command(); cmdName != "" {
		if cmd, ok := c.route[cmdName]; ok {
			msg.Text = cmd.Process(message.CommandArguments())
		} else {
			msg.Text = fmt.Sprintf("command [%s] not found", cmdName)
		}
	} else {
		msg.Text = fmt.Sprintf("your message - <%s>", message.Text)
	}
	_, err := c.bot.Send(msg)
	if err != nil {
		log.Printf("answer error: %v\n", err)
	}
}
