package receiver

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
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
	uid := uuid.New().String()
	ctx = helper.InjectUidPubToCtx(ctx, uid, in.GetPubSub().String())

	c.logger.Debugf("[%s] user create: [%s]", meta, in.User.String())

	user := adaptor.ToUserCoreModel(in.User).CreatedAtSet(time.Now().Unix())

	msg, err := json.Marshal(user)
	if err != nil {
		c.logger.Errorf("[%s] marshal err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserCreate),
		Value: sarama.ByteEncoder(msg),
	}); err != nil {
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserCreateResponse{
		Uid: uid,
	}, nil
}

func (c *core) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	meta := grpc.GetMetaFromContext(ctx)
	uid := uuid.New().String()
	ctx = helper.InjectUidPubToCtx(ctx, uid, in.GetPubSub().String())

	c.logger.Debugf("[%s] user create: [%s %s]", meta, in.GetName(), in.Profile.String())

	user := models.NewUser().
		NameSet(in.GetName()).
		PasswordSet(in.Profile.GetPassword()).
		EmailSet(in.Profile.GetEmail()).
		FullNameSet(in.Profile.GetFullName())

	msg, err := json.Marshal(user)
	if err != nil {
		c.logger.Errorf("[%s] marshal err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserUpdate),
		Value: sarama.ByteEncoder(msg),
	}); err != nil {
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserUpdateResponse{
		Uid: uid,
	}, nil
}

func (c *core) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	meta := grpc.GetMetaFromContext(ctx)
	uid := uuid.New().String()
	ctx = helper.InjectUidPubToCtx(ctx, uid, in.GetPubSub().String())

	c.logger.Debugf("[%s] user delete: [%s]", meta, in.GetName())

	if err := c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserDelete),
		Value: sarama.ByteEncoder(in.GetName()),
	}); err != nil {
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserDeleteResponse{
		Uid: uid,
	}, nil
}

func (c *core) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	meta := grpc.GetMetaFromContext(ctx)
	uid := uuid.New().String()
	ctx = helper.InjectUidPubToCtx(ctx, uid, in.GetPubSub().String())

	c.logger.Debugf("[%s] user get: [%s]", meta, in.GetName())

	if err := c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
		Topic: consts.TopicValidate,
		Key:   sarama.StringEncoder(consts.UserGet),
		Value: sarama.ByteEncoder(in.GetName()),
	}); err != nil {
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserGetResponse{
		Uid: uid,
	}, nil
}

func (c *core) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	meta := grpc.GetMetaFromContext(ctx)
	uid := uuid.New().String()
	ctx = helper.InjectUidPubToCtx(ctx, uid, in.GetPubSub().String())

	c.logger.Debugf("[%s] user list: [%v %v %v]", meta, in.GetLimit(), in.GetOffset(), in.GetOrder())

	params := models.NewUserListParams().
		LimitSet(in.GetLimit()).
		OffsetSet(in.GetOffset()).
		OrderSet(in.GetOrder())

	msg, err := json.Marshal(params)
	if err != nil {
		c.logger.Errorf("[%s] marshal err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
		Topic: consts.TopicData,
		Key:   sarama.StringEncoder(consts.UserList),
		Value: sarama.ByteEncoder(msg),
	}); err != nil {
		c.logger.Errorf("[%s] send message err: %v", meta, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UserListResponse{
		Uid: uid,
	}, nil
}

func (c *core) UserAllList(in *pb.UserAllListRequest, stream pb.User_UserAllListServer) error {
	meta := grpc.GetMetaFromContext(stream.Context())
	c.logger.Debugf("[%s] all users list: [%v %v]", meta, in.GetOrder(), in.GetLimit())

	dataStream, err := c.user.UserAllList(stream.Context(), &pb.UserAllListRequest{
		Order: in.GetOrder(),
		Limit: in.GetLimit(),
	})
	if err != nil {
		c.logger.Errorf("[%s] all user list: stream: %v", meta, err)
		return status.Error(codes.Internal, err.Error())
	}

	for {
		next, err := dataStream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			c.logger.Errorf("[%s] all users list: next chunk: %v", meta, err)
			return status.Error(codes.Internal, err.Error())
		}
		if err = stream.Send(next); err != nil {
			c.logger.Errorf("[%s] all users list: send chunk: %v", meta, err)
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func (c *core) sendMessageWithCtx(ctx context.Context, message *sarama.ProducerMessage) error {
	if err := helper.InjectHeaders(ctx, message); err != nil {
		return err
	}
	_, _, err := c.producer.SendMessage(message)
	return err
}
