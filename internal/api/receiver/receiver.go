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
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	pbModels "gitlab.ozon.dev/iTukaev/homework/pkg/api/models"
	"gitlab.ozon.dev/iTukaev/homework/pkg/grpc"
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
	meta := grpc.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user create: [%s]", meta, in.User.String())

	in.User.CreatedAt = time.Now().Unix()
	msg, err := json.Marshal(adaptor.ToUserCoreModel(in.User))
	if err != nil {
		c.logger.Errorf("[%s] marshal err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if _, _, err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserCreate),
		Value: sarama.ByteEncoder(msg),
	}); err != nil {
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return new(pb.UserCreateResponse), nil
}

func (c *core) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	meta := grpc.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user create: [%s %s]", meta, in.GetName(), in.Profile.String())

	msg, err := json.Marshal(adaptor.ToUserCoreModel(&pbModels.User{
		Name:     in.GetName(),
		Password: in.Profile.GetPassword(),
		Email:    in.Profile.GetEmail(),
		FullName: in.Profile.GetFullName(),
	}))
	if err != nil {
		c.logger.Errorf("[%s] marshal err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if _, _, err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserUpdate),
		Value: sarama.ByteEncoder(msg),
	}); err != nil {
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return new(pb.UserUpdateResponse), nil
}

func (c *core) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	meta := grpc.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user delete: [%s]", meta, in.GetName())

	if _, _, err := c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserDelete),
		Value: sarama.ByteEncoder(in.GetName()),
	}); err != nil {
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return new(pb.UserDeleteResponse), nil
}

func (c *core) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	meta := grpc.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user get: [%s]", meta, in.GetName())

	if in.GetName() == "" {
		c.logger.Errorf("[%s] empty [name]", meta)
		return nil, status.Error(codes.InvalidArgument, errors.New("field: [name] cannot be empty").Error())
	}

	resp, err := c.user.UserGet(ctx, in)
	if err != nil {
		c.logger.Errorf("[%s] user get: %v", meta, err)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	meta := grpc.GetMetaFromContext(ctx)
	c.logger.Debugf("[%s] user list: [%v %v %v]", meta, in.GetLimit(), in.GetOffset(), in.GetOrder())

	span := opentracing.StartSpan("receiver_user_list_" + meta)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	resp, err := c.user.UserList(ctx, in)
	if err != nil {
		c.logger.Errorf("[%s] user list: %v", meta, err)
		span.SetTag("error", true)
		return nil, err
	}

	return resp, nil
}

func (c *core) UserAllList(in *pb.UserAllListRequest, stream pb.User_UserAllListServer) error {
	meta := grpc.GetMetaFromContext(stream.Context())
	c.logger.Debugf("[%s] all users list: [%v %v]", meta, in.GetOrder(), in.GetLimit())

	span := opentracing.StartSpan("receiver_user_all_list_" + meta)
	ctx := opentracing.ContextWithSpan(stream.Context(), span)
	defer span.Finish()

	dataStream, err := c.user.UserAllList(ctx, &pb.UserAllListRequest{
		Order: in.GetOrder(),
		Limit: in.GetLimit(),
	})
	if err != nil {
		c.logger.Errorf("[%s] all user list: stream: %v", meta, err)
		span.SetTag("error", true)
		return status.Error(codes.Internal, err.Error())
	}

	for {
		next, err := dataStream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			c.logger.Errorf("[%s] all users list: next chunk: %v", meta, err)
			span.SetTag("error", true)
			return status.Error(codes.Internal, err.Error())
		}
		if err = stream.Send(next); err != nil {
			c.logger.Errorf("[%s] all users list: send chunk: %v", meta, err)
			span.SetTag("error", true)
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func (c *core) sendMessageWithCtx(
	ctx context.Context,
	message *sarama.ProducerMessage,
) (partition int32, offset int64, err error) {
	span := opentracing.SpanFromContext(ctx)
	if err = helper.InjectSpanIntoMessage(span, message); err != nil {
		return 0, 0, err
	}
	return c.producer.SendMessage(message)
}
