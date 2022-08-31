package get

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/status"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func New(api pb.UserClient, logger *zap.SugaredLogger) commandPkg.Interface {
	return &command{
		api:    api,
		logger: logger,
	}
}

type command struct {
	api    pb.UserClient
	logger *zap.SugaredLogger
}

func (c *command) Process(ctx context.Context, args string) string {
	args = strings.TrimSpace(args)
	if args == "" {
		return "invalid arguments"
	}

	user, err := c.api.UserGet(ctx, &pb.UserGetRequest{
		Name: args,
	})
	if err != nil {
		c.logger.Errorf("user [%s] get: %v\n", args, err)
		if st, ok := status.FromError(err); ok {
			return st.Message()
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
