package add

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	localPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/cache/local"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
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
	params := strings.Split(args, " ")
	if len(params) != 2 {
		return "invalid arguments"
	}

	if err := c.user.Create(ctx, models.User{
		Name:     params[0],
		Password: params[1],
	}); err != nil {
		log.Printf("user [%s] create: %v", params[0], err)
		if errors.Is(err, userPkg.ErrValidation) {
			return "invalid arguments"
		} else if errors.Is(err, localPkg.ErrUserAlreadyExists) {
			return err.Error()
		}
		return "internal error"
	}
	return fmt.Sprintf("user [%s] added", params[0])
}

func (*command) Name() string {
	return "add"
}

func (*command) Description() string {
	return "create user [/add <name> <password>]"
}
