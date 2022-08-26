package validator

import (
	"context"
	"io"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/iTukaev/homework/internal/counter"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

const (
	undefinedMeta = "undefined"

	userCreate  = "create"
	userUpdate  = "update"
	userDelete  = "delete"
	userGet     = "get"
	userList    = "list"
	userAllList = "all_list"
)

func New(user pb.UserClient, logger *zap.SugaredLogger) pb.UserServer {
	return &core{
		user:   user,
		logger: logger,
	}
}

type core struct {
	pr   sarama.AsyncProducer
	user pb.UserClient
	pb.UnimplementedUserServer
	logger *zap.SugaredLogger
}

func (c *core) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	counter.Request.Inc(userCreate)
	defer counter.Response.Inc(userCreate)

	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debugln(meta, "user create:", in.User.String())

	if in.User.GetName() == "" {
		counter.Errors.Inc(userCreate)
		c.logger.Errorln(meta, "empty [name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}
	if in.User.GetPassword() == "" {
		counter.Errors.Inc(userCreate)
		c.logger.Errorln(meta, "empty [password]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [password] cannot be empty").Error())
	}
	if in.User.GetEmail() == "" {
		counter.Errors.Inc(userCreate)
		c.logger.Errorln(meta, "empty [email]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [email] cannot be empty").Error())
	}
	if in.User.GetFullName() == "" {
		counter.Errors.Inc(userCreate)
		c.logger.Errorln(meta, "empty [full_name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [full_name] cannot be empty").Error())
	}
	in.User.CreatedAt = time.Now().Unix()

	resp, err := c.user.UserCreate(ctx, in)
	if err != nil {
		counter.Errors.Inc(userCreate)
		c.logger.Errorln(meta, "user create:", err)
		return nil, err
	}

	counter.Success.Inc(userCreate)
	return resp, nil
}

func (c *core) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	counter.Request.Inc(userUpdate)
	defer counter.Response.Inc(userUpdate)

	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debugln(meta, "user update:", in.GetName(), in.Profile.String())

	if in.GetName() == "" {
		counter.Errors.Inc(userUpdate)
		c.logger.Errorln(meta, "empty [name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserUpdate(ctx, in)
	if err != nil {
		counter.Errors.Inc(userUpdate)
		c.logger.Errorln(meta, "user update:", err)
		return nil, err
	}

	counter.Success.Inc(userUpdate)
	return resp, nil
}

func (c *core) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	counter.Request.Inc(userDelete)
	defer counter.Response.Inc(userDelete)

	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debugln(meta, "user delete:", in.GetName())

	if in.GetName() == "" {
		counter.Errors.Inc(userDelete)
		c.logger.Errorln(meta, "empty [name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserDelete(ctx, in)
	if err != nil {
		counter.Errors.Inc(userDelete)
		c.logger.Errorln(meta, "user delete:", err)
		return nil, err
	}

	counter.Success.Inc(userDelete)
	return resp, nil
}

func (c *core) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	counter.Request.Inc(userGet)
	defer counter.Response.Inc(userGet)

	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debugln(meta, "user get:", in.GetName())

	if in.GetName() == "" {
		counter.Errors.Inc(userDelete)
		c.logger.Errorln(meta, "empty [name]:")
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserGet(ctx, in)
	if err != nil {
		counter.Errors.Inc(userDelete)
		c.logger.Errorln(meta, "user get:", err)
		return nil, err
	}

	counter.Success.Inc(userDelete)
	return resp, nil
}

func (c *core) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	counter.Request.Inc(userList)
	defer counter.Response.Inc(userList)

	meta, ok := ctx.Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debugln(meta, "user get:", in.GetLimit(), in.GetOffset(), in.GetOrder())

	resp, err := c.user.UserList(ctx, in)
	if err != nil {
		counter.Errors.Inc(userList)
		c.logger.Errorln(meta, "user list:", err)
		return nil, err
	}

	counter.Success.Inc(userList)
	return resp, nil
}

func (c *core) UserAllList(in *pb.UserAllListRequest, stream pb.User_UserAllListServer) error {
	counter.Request.Inc(userAllList)
	defer counter.Response.Inc(userAllList)

	meta, ok := stream.Context().Value("meta").(string)
	if !ok {
		meta = undefinedMeta
	}
	c.logger.Debugln(meta, "all users list", in.GetOrder(), in.GetLimit())

	dataStream, err := c.user.UserAllList(stream.Context(), &pb.UserAllListRequest{
		Order: in.GetOrder(),
		Limit: in.GetLimit(),
	})
	if err != nil {
		counter.Errors.Inc(userAllList)
		c.logger.Errorln(meta, "all users list, stream", err)
		return status.Error(codes.Internal, err.Error())
	}

	for {
		next, err := dataStream.Recv()
		if errors.Is(err, io.EOF) {
			counter.Success.Inc(userAllList)
			return nil
		}
		if err != nil {
			counter.Errors.Inc(userAllList)
			c.logger.Errorln(meta, "all users list, next chunk", err)
			return err
		}
		if err = stream.Send(next); err != nil {
			counter.Errors.Inc(userAllList)
			c.logger.Errorln(meta, "all users list, send chunk", err)
			return status.Error(codes.Internal, err.Error())
		}
	}
}
