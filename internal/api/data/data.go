package data

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	grpcPkg "gitlab.ozon.dev/iTukaev/homework/pkg/grpc"
)

func New(user userPkg.Interface, logger *zap.SugaredLogger) pb.UserServer {
	return &core{
		user:   user,
		logger: logger,
	}
}

type core struct {
	user   userPkg.Interface
	logger *zap.SugaredLogger
	pb.UnimplementedUserServer
}

func (c *core) UserAllList(in *pb.UserAllListRequest, stream pb.User_UserAllListServer) error {
	meta := grpcPkg.GetMetaFromContext(stream.Context())
	c.logger.Debugln(meta, "all users list", in.GetOrder(), in.GetLimit())

	offset := uint64(0)
	for {
		users, err := c.user.List(stream.Context(), in.GetOrder(), in.GetLimit(), offset)
		if err != nil {
			c.logger.Errorln(meta, "get list", err)
			return status.Error(codes.Internal, err.Error())
		}

		if len(users) == 0 {
			return nil
		}

		if err = stream.Send(&pb.UserAllListResponse{
			Users: adaptor.ToUserListPbModel(users),
		}); err != nil {
			c.logger.Errorln(meta, "all users list, send chunk", err)
			return status.Error(codes.Internal, err.Error())
		}
		offset++
	}
}
