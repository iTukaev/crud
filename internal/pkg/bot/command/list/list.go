package list

import (
	"context"
	"strconv"
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
	params := strings.Split(args, " ")
	if len(params) != 3 {
		return "invalid arguments"
	}
	order, err := strconv.ParseBool(params[0])
	if err != nil {
		return "invalid [order] argument"
	}
	limit, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return "invalid [limit] argument"
	}
	offset, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		return "invalid [offset] argument"
	}

	list, err := c.api.UserList(ctx, &pb.UserListRequest{
		Order:  order,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.logger.Errorf("user list, arguments [%s]: %v\n", args, err)
		if st, ok := status.FromError(err); ok {
			return st.Message()
		}
		return "internal error"
	}

	//for i, u := range list.Users {
	//	result += u.String()
	//	if i != len(list.Users)-1 {
	//		result += "\n"
	//	}
	//}
	return list.GetUid()
}

func (*command) Name() string {
	return "list"
}

func (*command) Description() string {
	return "get all users info [/list <order-true/false> <limit> <offset>]"
}
