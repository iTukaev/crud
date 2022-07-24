package get

import (
	"log"

	"github.com/pkg/errors"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	localPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/cache/local"
)

func New(user userPkg.Interface) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.Interface
}

func (c *command) Process(args string) string {
	user, err := c.user.Get(args)
	if err != nil {
		if errors.Is(err, userPkg.ErrValidation) {
			return "invalid arguments"
		} else if errors.Is(err, localPkg.ErrUserNotExists) || errors.Is(err, localPkg.ErrUserExists) {
			return err.Error()
		}
		log.Printf("user [%s] get: %v", args, err)
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
