package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
)

const (
	contextTimeout = 5 * time.Second
)

type Interface interface {
	RegisterCommand(cmd commandPkg.Interface)
	Run(ctx context.Context)
	Stop()
}

func MustNew(id string) Interface {
	bot, err := tgbotapi.NewBotAPI(id)
	if err != nil {
		log.Panic(errors.Wrap(err, "new API bot"))
	}

	bot.Debug = false

	return &commander{
		bot:   bot,
		route: make(map[string]commandPkg.Interface),
	}
}

type commander struct {
	bot   *tgbotapi.BotAPI
	route map[string]commandPkg.Interface
}

// RegisterCommand - not thread safe
func (c *commander) RegisterCommand(cmd commandPkg.Interface) {
	c.route[cmd.Name()] = cmd
}

func (c *commander) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go c.handleMessage(ctx, update.Message)
		}
	}
}

func (c *commander) Stop() {
	c.bot.StopReceivingUpdates()
}

func (c *commander) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	if cmdName := message.Command(); cmdName != "" {
		if cmd, ok := c.route[cmdName]; ok {
			msg.Text = cmd.Process(ctxWithTimeout, message.CommandArguments())
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
