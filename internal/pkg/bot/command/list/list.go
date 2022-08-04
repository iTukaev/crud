package list

import (
	"context"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
)

func New(user userPkg.Interface) commandPkg.Interface {
	return &command{
		user: user,
	}
}

type command struct {
	user userPkg.Interface
}

func (c *command) Process(ctx context.Context, _ string) string {
	list, err := c.user.List(ctx, false, 10, 1)
	if err != nil {
		return err.Error()
	}

	result := ""
	for i, u := range list {
		result += u.String()
		if i != len(list)-1 {
			result += "\n"
		}
	}
	return result
}

func (*command) Name() string {
	return "list"
}

func (*command) Description() string {
	return "get all users info"
}
