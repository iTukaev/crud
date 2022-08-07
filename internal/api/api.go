package api

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
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	pbModels "gitlab.ozon.dev/iTukaev/homework/pkg/api/models"
)

const (
	contextTimeout = 5 * time.Second
)

func New(user userPkg.Interface) pb.UserServer {
	return &implementation{
		user: user,
	}
}

type implementation struct {
	user userPkg.Interface
	pb.UnimplementedUserServer
}

func (i *implementation) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if err := i.user.Create(ctx, models.User{
		Name:     in.User.GetName(),
		Password: in.User.GetPassword(),
		Email:    in.User.GetEmail(),
		FullName: in.User.GetFullName(),
	}); err != nil {
		log.Printf("user [%s] create: %v", in.User.GetName(), err)

		switch {
		case errors.Is(err, userPkg.ErrValidation):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errorsPkg.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, errorsPkg.ErrTimeout):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserCreateResponse{}, nil
}

func (i *implementation) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if err := i.user.Update(ctx, models.User{
		Name:     in.GetName(),
		Password: in.Profile.GetPassword(),
		Email:    in.Profile.GetEmail(),
		FullName: in.Profile.GetFullName(),
	}); err != nil {
		log.Printf("user [%s] update: %v", in.GetName(), err)

		switch {
		case errors.Is(err, userPkg.ErrValidation), errors.Is(err, errorsPkg.ErrUserNotFound):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errorsPkg.ErrTimeout):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserUpdateResponse{}, nil
}

func (i *implementation) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	if err := i.user.Delete(ctx, in.GetName()); err != nil {
		log.Printf("user [%s] delete: %v", in.GetName(), err)

		switch {
		case errors.Is(err, userPkg.ErrValidation), errors.Is(err, errorsPkg.ErrUserNotFound):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errorsPkg.ErrTimeout):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserDeleteResponse{}, nil
}

func (i *implementation) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	user, err := i.user.Get(ctx, in.GetName())
	if err != nil {
		log.Printf("user [%s] get: %v", in.GetName(), err)

		switch {
		case errors.Is(err, userPkg.ErrValidation), errors.Is(err, errorsPkg.ErrUserNotFound):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errorsPkg.ErrTimeout):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserGetResponse{
		User: &pbModels.User{
			Name:      user.Name,
			Password:  user.Password,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (i *implementation) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	users, err := i.user.List(ctx, in.GetOrder(), in.GetLimit(), in.GetOffset())
	if err != nil {
		log.Printf("user list: %v", err)
		if errors.Is(err, errorsPkg.ErrTimeout) {
			return &pb.UserListResponse{}, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return &pb.UserListResponse{}, status.Error(codes.Internal, err.Error())
	}

	resp := make([]*pbModels.User, 0, len(users))
	for _, user := range users {
		resp = append(resp, &pbModels.User{
			Name:      user.Name,
			Password:  user.Password,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
		})
	}

	return &pb.UserListResponse{
		Users: resp,
	}, nil
}
