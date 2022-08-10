package add

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc/status"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	pbModels "gitlab.ozon.dev/iTukaev/homework/pkg/api/models"
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
	params := strings.Split(args, " ")
	if len(params) != 4 {
		return "invalid arguments"
	}

	if _, err := c.api.UserCreate(ctx, &pb.UserCreateRequest{
		User: &pbModels.User{
			Name:     params[0],
			Password: params[1],
			Email:    params[2],
			FullName: params[3],
		},
	}); err != nil {
		log.Printf("user [%s] create: %v", params[0], err)
		if st, ok := status.FromError(err); ok {
			return st.Message()
		}
		return "internal error"
	}
	return fmt.Sprintf("user [%s] added", params[0])
}

func (*command) Name() string {
	return "add"
}

func (*command) Description() string {
	return "create user [/add <name> <password> <email> <full_name>]"
}
