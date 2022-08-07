package update

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
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
	params := strings.Split(args, " ")
	if len(params) != 2 {
		return "invalid arguments"
	}

	if err := c.user.Update(ctx, models.User{
		Name:     params[0],
		Password: params[1],
	}); err != nil {
		log.Printf("user [%s] update: %v", params[0], err)
		if errors.Is(err, userPkg.ErrValidation) {
			return "invalid arguments"
		} else if errors.Is(err, errorsPkg.ErrUserNotFound) {
			return err.Error()
		}
		return "internal error"
	}
	return fmt.Sprintf("user [%s] updated", params[0])
}

func (*command) Name() string {
	return "update"
}

func (*command) Description() string {
	return "update user [/update <new name> <new password>]"
}
