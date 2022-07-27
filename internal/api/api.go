package api

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	localPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/cache/local"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
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
	if err := i.user.Create(models.User{
		Name:     in.GetName(),
		Password: in.GetPassword(),
	}); err != nil {
		log.Printf("user [%s] create: %v", in.GetName(), err)
		switch {
		case errors.Is(err, userPkg.ErrValidation):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, localPkg.ErrUserAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserCreateResponse{}, nil
}

func (i *implementation) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	if err := i.user.Update(models.User{
		Name:     in.GetName(),
		Password: in.GetPassword(),
	}); err != nil {
		log.Printf("user [%s] update: %v", in.GetName(), err)
		switch {
		case errors.Is(err, userPkg.ErrValidation):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, localPkg.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserUpdateResponse{}, nil
}

func (i *implementation) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	if err := i.user.Delete(in.GetName()); err != nil {
		log.Printf("user [%s] delete: %v", in.GetName(), err)
		switch {
		case errors.Is(err, userPkg.ErrValidation):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, localPkg.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserDeleteResponse{}, nil
}

func (i *implementation) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	user, err := i.user.Get(in.GetName())
	if err != nil {
		log.Printf("user [%s] get: %v", in.GetName(), err)
		switch {
		case errors.Is(err, userPkg.ErrValidation):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, localPkg.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserGetResponse{
		User: &pb.UserGetResponse_User{
			Name:     user.Name,
			Password: user.Password,
		},
	}, nil
}

func (i *implementation) UserList(ctx context.Context, _ *pb.UserListRequest) (*pb.UserListResponse, error) {
	users := i.user.List()

	resp := make([]*pb.UserListResponse_User, 0, len(users))
	for _, user := range users {
		resp = append(resp, &pb.UserListResponse_User{
			Name:     user.Name,
			Password: user.Password,
		})
	}
	return &pb.UserListResponse{
		Users: resp,
	}, nil
}
