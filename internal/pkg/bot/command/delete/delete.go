package delete

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc/status"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func New(api pb.UserClient) commandPkg.Interface {
	return &command{
		api: api,
	}
}

type command struct {
	api pb.UserClient
}

func (c *command) Process(ctx context.Context, args string) string {
	args = strings.TrimSpace(args)
	if args == "" {
		return "invalid arguments"
	}

	if _, err := c.api.UserDelete(ctx, &pb.UserDeleteRequest{
		Name: args,
	}); err != nil {
		log.Printf("user [%s] delete: %v\n", args, err)
		if st, ok := status.FromError(err); ok {
			return st.Message()
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
