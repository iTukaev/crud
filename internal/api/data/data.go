package data

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/customerrors"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

const (
	contextTimeout = 5 * time.Second

	dataService = "data"
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

func (c *core) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	ctx, meta := helper.GetMetaFromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	go helper.StartNewSpan(ctx, dataService, stop)

	user, err := c.user.Get(ctx, in.GetName())
	if err != nil {
		c.logger.Errorln(meta, "user get:", err)

		switch {
		case errors.Is(err, errorsPkg.ErrUserNotFound):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errorsPkg.ErrTimeout):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserGetResponse{
		User: adaptor.ToUserPbModel(user),
	}, nil
}

func (c *core) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	ctx, meta := helper.GetMetaFromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	c.logger.Debugln(meta, "user get:", in.GetLimit(), in.GetOffset(), in.GetOrder())

	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	go helper.StartNewSpan(ctx, dataService, stop)

	users, err := c.user.List(ctx, in.GetOrder(), in.GetLimit(), in.GetOffset())
	if err != nil {
		c.logger.Errorln(meta, "user list:", err)
		if errors.Is(err, errorsPkg.ErrTimeout) {
			return &pb.UserListResponse{}, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return &pb.UserListResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserListResponse{
		Users: adaptor.ToUserListPbModel(users),
	}, nil
}

func (c *core) UserAllList(in *pb.UserAllListRequest, stream pb.User_UserAllListServer) error {
	ctx, meta := helper.GetMetaFromContext(stream.Context())
	c.logger.Debugln(meta, "all users list", in.GetOrder(), in.GetLimit())

	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	go helper.StartNewSpan(ctx, dataService, stop)

	offset := uint64(0)
	for {
		users, err := c.user.List(ctx, in.GetOrder(), in.GetLimit(), offset)
		if err != nil {
			c.logger.Errorln(meta, "Get list", err)
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
