package validator

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func New(user pb.UserClient) pb.UserServer {
	return &core{
		user: user,
	}
}

type core struct {
	user pb.UserClient
	pb.UnimplementedUserServer
}

func (c *core) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	if in.User.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}
	if in.User.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [password] cannot be empty").Error())
	}
	if in.User.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [email] cannot be empty").Error())
	}
	if in.User.GetFullName() == "" {
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [full_name] cannot be empty").Error())
	}
	in.User.CreatedAt = time.Now().Unix()

	resp, err := c.user.UserCreate(ctx, in)
	if err != nil {
		log.Printf("user [%s] create: %v\n", in.User.GetName(), err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserUpdate(ctx, in)
	if err != nil {
		log.Printf("user [%s] update: %v\n", in.GetName(), err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserDelete(ctx, in)
	if err != nil {
		log.Printf("user [%s] delete: %v\n", in.GetName(), err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserGet(ctx, in)
	if err != nil {
		log.Printf("user [%s] get: %v\n", in.GetName(), err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	resp, err := c.user.UserList(ctx, in)
	if err != nil {
		log.Printf("user list: %v\n", err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserAllList(in *pb.UserAllListRequest, stream pb.User_UserAllListServer) error {
	dataStream, err := c.user.UserAllList(stream.Context(), &pb.UserAllListRequest{
		Order: in.GetOrder(),
		Limit: in.GetLimit(),
	})
	if err != nil {
		log.Printf("all users list, stream: %v\n", err)
		return status.Error(codes.Internal, err.Error())
	}

	for {
		next, err := dataStream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			log.Printf("all users list, next chunk: %v\n", err)
			return err
		}
		if err = stream.Send(next); err != nil {
			log.Printf("all users list, send chunk: %v\n", err)
			return status.Error(codes.Internal, err.Error())
		}
	}
}
