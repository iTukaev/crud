package delete

import (
	"fmt"
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
	if err := c.user.Delete(args); err != nil {
		log.Printf("user [%s] delete: %v", args, err)
		if errors.Is(err, userPkg.ErrValidation) {
			return "invalid arguments"
		} else if errors.Is(err, localPkg.ErrUserNotFound) {
			return err.Error()
		}
		return "internal error"
	}
	return fmt.Sprintf("user [%s] deleted", args)
}

func (*command) Name() string {
	return "del"
}

func (*command) Description() string {
	return "delete user [/del <name>]"
}
