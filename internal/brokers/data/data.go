package data

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/customerrors"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

const (
	brokerDataService = "broker_data"
)

type sender interface {
	userCreate(ctx context.Context, msg *sarama.ConsumerMessage) error
	userUpdate(ctx context.Context, msg *sarama.ConsumerMessage) error
	userDelete(ctx context.Context, msg *sarama.ConsumerMessage) error
	userGet(ctx context.Context, msg *sarama.ConsumerMessage) error
	userList(ctx context.Context, msg *sarama.ConsumerMessage) error
}

func newSender(user userPkg.Interface, logger *zap.SugaredLogger, producer sarama.SyncProducer) sender {
	return &core{
		user:     user,
		producer: producer,
		logger:   logger,
	}
}

type core struct {
	user     userPkg.Interface
	producer sarama.SyncProducer
	logger   *zap.SugaredLogger
}

func (c *core) userCreate(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, brokerDataService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	c.logger.Debugf("user [%s]", user.String())

	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.StringEncoder(consts.UserCreate),
	}

	if err := c.user.Create(ctx, user); err != nil {
		if errors.Is(err, errorsPkg.ErrUserAlreadyExists) {
			c.logger.Errorf("user create: %v", err)
			return c.sendErrorWithCtx(ctx, message, err.Error())
		}
		return err
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) userUpdate(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, brokerDataService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "message unmarshal")
	}

	c.logger.Debugf("user [%s]", user.String())

	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.StringEncoder(consts.UserUpdate),
	}

	if err := c.user.Update(ctx, user); err != nil {
		if errors.Is(err, errorsPkg.ErrUserNotFound) {
			c.logger.Errorf("user update: %v", err)
			return c.sendErrorWithCtx(ctx, message, err.Error())
		}
		return err
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) userDelete(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, brokerDataService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	name := string(msg.Value)

	c.logger.Debugf("name: [%s]", name)

	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.StringEncoder(consts.UserDelete),
	}

	if err := c.user.Delete(ctx, name); err != nil {
		if errors.Is(err, errorsPkg.ErrUserNotFound) {
			c.logger.Errorf("user delete: %v", err)
			return c.sendErrorWithCtx(ctx, message, err.Error())
		}
		return err
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) userGet(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, brokerDataService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	name := string(msg.Value)

	c.logger.Debugf("name: [%s]", name)

	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.StringEncoder(consts.UserGet),
	}

	user, err := c.user.Get(ctx, name)
	if err != nil {
		if errors.Is(err, errorsPkg.ErrUserNotFound) {
			c.logger.Errorf("user get: %v", err)
			return c.sendErrorWithCtx(ctx, message, err.Error())
		}
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "marshal user")
	}
	message.Value = sarama.ByteEncoder(data)

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) userList(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, brokerDataService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	params := models.NewUserListParams()
	if err := json.Unmarshal(msg.Value, params); err != nil {
		return errors.Wrap(err, "unmarshal list parameters")
	}

	c.logger.Debugf("parameters: [%d %d %v]", params.Limit, params.Offset, params.Order)

	list, err := c.user.List(ctx, params.Order, params.Limit, params.Offset)
	if err != nil {
		return err
	}

	data, err := json.Marshal(list)
	if err != nil {
		return errors.Wrap(err, "marshal user")
	}
	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.StringEncoder(consts.UserList),
		Value: sarama.ByteEncoder(data),
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) sendErrorWithCtx(
	ctx context.Context,
	message *sarama.ProducerMessage,
	description string,
) error {
	if err := helper.InjectHeaders(ctx, message); err != nil {
		return err
	}
	message.Topic = consts.TopicError
	message.Value = sarama.StringEncoder(description)

	_, _, err := c.producer.SendMessage(message)
	return err
}

func (c *core) sendMessageWithCtx(ctx context.Context, message *sarama.ProducerMessage) error {
	if err := helper.InjectHeaders(ctx, message); err != nil {
		return err
	}
	_, _, err := c.producer.SendMessage(message)
	return err
}
