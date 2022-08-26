package delete

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/status"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

const (
	deleteName        = "del"
	deleteDescription = "delete user [/del <name>]"
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

	if _, err := c.api.UserDelete(ctx, &pb.UserDeleteRequest{
		Name: args,
	}); err != nil {
		c.logger.Errorf("user [%s] delete: %v\n", args, err)
		if st, ok := status.FromError(err); ok {
			return st.Message()
		}
		return "internal error"
	}
	return fmt.Sprintf("user [%s] deleted", args)
}

func (*command) Name() string {
	return deleteName
}

func (*command) Description() string {
	return deleteDescription
}
