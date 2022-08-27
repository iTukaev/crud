package receiver

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	"gitlab.ozon.dev/iTukaev/homework/internal/counter"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	pbModels "gitlab.ozon.dev/iTukaev/homework/pkg/api/models"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

func New(user pb.UserClient, logger *zap.SugaredLogger, producer sarama.SyncProducer) pb.UserServer {
	return &core{
		producer: producer,
		user:     user,
		logger:   logger,
	}
}

type core struct {
	producer sarama.SyncProducer
	user     pb.UserClient
	pb.UnimplementedUserServer
	logger *zap.SugaredLogger
}

func (c *core) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	counter.Request.Inc(consts.UserCreate)
	defer counter.Response.Inc(consts.UserCreate)

	_, meta := helper.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user create: [%s]", meta, in.User.String())

	span := opentracing.StartSpan("receiver_user_create_" + meta)
	defer span.Finish()

	in.User.CreatedAt = time.Now().Unix()
	msg, err := json.Marshal(adaptor.ToUserCoreModel(in.User))
	if err != nil {
		counter.Errors.Inc(consts.UserCreate)
		c.logger.Errorf("[%s] marshal err: %v", meta, err)
		span.SetTag("error", true)
		return nil, status.Error(codes.Internal, err.Error())
	}

	message := &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserCreate),
		Value: sarama.ByteEncoder(msg),
	}
	if err = helper.InjectSpanIntoMessage(span, message); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if _, _, err = c.producer.SendMessage(message); err != nil {
		counter.Errors.Inc(consts.UserCreate)
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		span.SetTag("error", true)
		return nil, status.Error(codes.Internal, err.Error())
	}

	counter.Success.Inc(consts.UserCreate)
	return new(pb.UserCreateResponse), nil
}

func (c *core) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	counter.Request.Inc(consts.UserUpdate)
	defer counter.Response.Inc(consts.UserUpdate)

	_, meta := helper.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user create: [%s %s]", meta, in.GetName(), in.Profile.String())

	span := opentracing.StartSpan("receiver_user_update_" + meta)
	defer span.Finish()

	msg, err := json.Marshal(adaptor.ToUserCoreModel(&pbModels.User{
		Name:     in.GetName(),
		Password: in.Profile.GetPassword(),
		Email:    in.Profile.GetEmail(),
		FullName: in.Profile.GetFullName(),
	}))
	if err != nil {
		counter.Errors.Inc(consts.UserUpdate)
		c.logger.Errorf("[%s] marshal err: %v", meta, err)
		span.SetTag("error", true)
		return nil, status.Error(codes.Internal, err.Error())
	}

	message := &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserUpdate),
		Value: sarama.ByteEncoder(msg),
	}
	if err = helper.InjectSpanIntoMessage(span, message); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, _, err = c.producer.SendMessage(message)
	if err != nil {
		counter.Errors.Inc(consts.UserUpdate)
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		span.SetTag("error", true)
		return nil, status.Error(codes.Internal, err.Error())
	}

	counter.Success.Inc(consts.UserUpdate)
	return new(pb.UserUpdateResponse), nil
}

func (c *core) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	counter.Request.Inc(consts.UserDelete)
	defer counter.Response.Inc(consts.UserDelete)

	_, meta := helper.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user delete: [%s]", meta, in.GetName())

	span := opentracing.StartSpan("receiver_user_delete_" + meta)
	defer span.Finish()

	message := &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserDelete),
		Value: sarama.ByteEncoder(in.GetName()),
	}
	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, _, err := c.producer.SendMessage(message)
	if err != nil {
		counter.Errors.Inc(consts.UserDelete)
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		span.SetTag("error", true)
		return nil, status.Error(codes.Internal, err.Error())
	}

	counter.Success.Inc(consts.UserDelete)
	return new(pb.UserDeleteResponse), nil
}

func (c *core) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	counter.Request.Inc(consts.UserGet)
	defer counter.Response.Inc(consts.UserGet)

	ctx, meta := helper.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user get: [%s]", meta, in.GetName())

	span := opentracing.StartSpan("receiver_user_get_" + meta)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	if in.GetName() == "" {
		counter.Errors.Inc(consts.UserGet)
		c.logger.Errorf("[%s] empty [name]", meta)
		span.SetTag("error", true)
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserGet(ctx, in)
	if err != nil {
		counter.Errors.Inc(consts.UserGet)
		c.logger.Errorf("[%s] user get: %v", meta, err)
		span.SetTag("error", true)
		return nil, err
	}

	counter.Success.Inc(consts.UserGet)
	return resp, nil
}

func (c *core) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	counter.Request.Inc(consts.UserList)
	defer counter.Response.Inc(consts.UserList)

	ctx, meta := helper.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user list: [%v %v %v]", meta, in.GetLimit(), in.GetOffset(), in.GetOrder())

	span := opentracing.StartSpan("receiver_user_list_" + meta)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	resp, err := c.user.UserList(ctx, in)
	if err != nil {
		counter.Errors.Inc(consts.UserList)
		c.logger.Errorf("[%s] user list: %v", meta, err)
		span.SetTag("error", true)
		return nil, err
	}

	counter.Success.Inc(consts.UserList)
	return resp, nil
}

func (c *core) UserAllList(in *pb.UserAllListRequest, stream pb.User_UserAllListServer) error {
	counter.Request.Inc(consts.UserAllList)
	defer counter.Response.Inc(consts.UserAllList)

	ctx, meta := helper.GetMetaFromContext(stream.Context())
	c.logger.Debugf("[%s] all users list: [%v %v]", meta, in.GetOrder(), in.GetLimit())

	span := opentracing.StartSpan("receiver_user_all_list_" + meta)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	dataStream, err := c.user.UserAllList(ctx, &pb.UserAllListRequest{
		Order: in.GetOrder(),
		Limit: in.GetLimit(),
	})
	if err != nil {
		counter.Errors.Inc(consts.UserAllList)
		c.logger.Errorf("[%s] all user list: stream: %v", meta, err)
		span.SetTag("error", true)
		return status.Error(codes.Internal, err.Error())
	}

	for {
		next, err := dataStream.Recv()
		if errors.Is(err, io.EOF) {
			counter.Success.Inc(consts.UserAllList)
			return nil
		}
		if err != nil {
			counter.Errors.Inc(consts.UserAllList)
			c.logger.Errorf("[%s] all users list: next chunk: %v", meta, err)
			span.SetTag("error", true)
			return status.Error(codes.Internal, err.Error())
		}
		if err = stream.Send(next); err != nil {
			counter.Errors.Inc(consts.UserAllList)
			c.logger.Errorf("[%s] all users list: send chunk: %v", meta, err)
			span.SetTag("error", true)
			return status.Error(codes.Internal, err.Error())
		}
	}
}
