package data

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

const (
	contextTimeout = 5 * time.Second
)

func New(user userPkg.Interface) pb.UserServer {
	return &core{
		user: user,
	}
}

type core struct {
	user userPkg.Interface
	pb.UnimplementedUserServer
}

func (c *core) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if err := c.user.Create(ctx, models.User{
		Name:      in.User.GetName(),
		Password:  in.User.GetPassword(),
		Email:     in.User.GetEmail(),
		FullName:  in.User.GetFullName(),
		CreatedAt: in.User.GetCreatedAt(),
	}); err != nil {
		log.Printf("user [%s] create: %v\n", in.User.GetName(), err)

		switch {
		case errors.Is(err, errorsPkg.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, errorsPkg.ErrTimeout):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserCreateResponse{}, nil
}

func (c *core) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if err := c.user.Update(ctx, models.User{
		Name:     in.GetName(),
		Password: in.Profile.GetPassword(),
		Email:    in.Profile.GetEmail(),
		FullName: in.Profile.GetFullName(),
	}); err != nil {
		log.Printf("user [%s] update: %v\n", in.GetName(), err)

		switch {
		case errors.Is(err, errorsPkg.ErrUserNotFound):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errorsPkg.ErrTimeout):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserUpdateResponse{}, nil
}

func (c *core) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if err := c.user.Delete(ctx, in.GetName()); err != nil {
		log.Printf("user [%s] delete: %v\n", in.GetName(), err)

		switch {
		case errors.Is(err, errorsPkg.ErrUserNotFound):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errorsPkg.ErrTimeout):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserDeleteResponse{}, nil
}

func (c *core) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	user, err := c.user.Get(ctx, in.GetName())
	if err != nil {
		log.Printf("user [%s] get: %v\n", in.GetName(), err)

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
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	users, err := c.user.List(ctx, in.GetOrder(), in.GetLimit(), in.GetOffset())
	if err != nil {
		log.Printf("user list: %v\n", err)
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
	offset := uint64(0)
	for {
		users, err := c.user.List(stream.Context(), in.GetOrder(), in.GetLimit(), offset)
		if err != nil {
			log.Printf("user list: %v\n", err)
			return status.Error(codes.Internal, err.Error())
		}

		if len(users) == 0 {
			return nil
		}

		if err = stream.Send(&pb.UserAllListResponse{
			Users: adaptor.ToUserListPbModel(users),
		}); err != nil {
			log.Printf("stream send user list: %v\n", err)
			return status.Error(codes.Internal, err.Error())
		}
		offset++
	}
}
