package validator

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

const (
	undefinedMeta = "undefined"
)

func New(user pb.UserClient, logger loggerPkg.Interface) pb.UserServer {
	return &core{
		user:   user,
		logger: logger,
	}
}

type core struct {
	user pb.UserClient
	pb.UnimplementedUserServer
	logger loggerPkg.Interface
}

func (c *core) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debug(meta, "user create:", in.User.String())

	if in.User.GetName() == "" {
		c.logger.Error(meta, "empty [name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}
	if in.User.GetPassword() == "" {
		c.logger.Error(meta, "empty [password]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [password] cannot be empty").Error())
	}
	if in.User.GetEmail() == "" {
		c.logger.Error(meta, "empty [email]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [email] cannot be empty").Error())
	}
	if in.User.GetFullName() == "" {
		c.logger.Error(meta, "empty [full_name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [full_name] cannot be empty").Error())
	}
	in.User.CreatedAt = time.Now().Unix()

	resp, err := c.user.UserCreate(ctx, in)
	if err != nil {
		c.logger.Error(meta, "user create:", err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debug(meta, "user update:", in.GetName(), in.Profile.String())

	if in.GetName() == "" {
		c.logger.Error(meta, "empty [name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserUpdate(ctx, in)
	if err != nil {
		c.logger.Error(meta, "user update:", err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debug(meta, "user delete:", in.GetName())

	if in.GetName() == "" {
		c.logger.Error(meta, "empty [name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserDelete(ctx, in)
	if err != nil {
		c.logger.Error(meta, "user delete:", err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debug(meta, "user get:", in.GetName())

	if in.GetName() == "" {
		c.logger.Error(meta, "empty [name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserGet(ctx, in)
	if err != nil {
		c.logger.Error(meta, "user get:", err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debug(meta, "user get:", in.GetLimit(), in.GetOffset(), in.GetOrder())

	resp, err := c.user.UserList(ctx, in)
	if err != nil {
		c.logger.Error(meta, "user list:", err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserAllList(in *pb.UserAllListRequest, stream pb.User_UserAllListServer) error {
	meta, ok := stream.Context().Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debug(meta, "all users list", in.GetOrder(), in.GetLimit())

	dataStream, err := c.user.UserAllList(stream.Context(), &pb.UserAllListRequest{
		Order: in.GetOrder(),
		Limit: in.GetLimit(),
	})
	if err != nil {
		c.logger.Error(meta, "all users list, stream", err)
		return status.Error(codes.Internal, err.Error())
	}

	for {
		next, err := dataStream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			c.logger.Error(meta, "all users list, next chunk", err)
			return err
		}
		if err = stream.Send(next); err != nil {
			c.logger.Error(meta, "all users list, send chunk", err)
			return status.Error(codes.Internal, err.Error())
		}
	}
}
