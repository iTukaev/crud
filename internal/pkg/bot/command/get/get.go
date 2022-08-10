package get

import (
	"context"
	"log"

	"github.com/pkg/errors"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
)

func New(user userPkg.Interface) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.Interface
}

func (c *command) Process(ctx context.Context, args string) string {
	user, err := c.user.Get(ctx, args)
	if err != nil {
		log.Printf("user [%s] get: %v", args, err)
		if errors.Is(err, userPkg.ErrValidation) {
			return "invalid arguments"
		} else if errors.Is(err, errorsPkg.ErrUserNotFound) {
			return err.Error()
		}
		return "internal error"
	}
	return user.String()
}

func (*command) Name() string {
	return "get"
}

func (*command) Description() string {
	return "get user info [/get <name>]"
}
