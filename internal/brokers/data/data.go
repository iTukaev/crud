package data

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

const (
	brokerDataService = "broker_data"
)

type sender interface {
	userCreate(msg *sarama.ConsumerMessage) error
	userUpdate(msg *sarama.ConsumerMessage) error
	userDelete(msg *sarama.ConsumerMessage) error
}

func newSender(ctx context.Context, user userPkg.Interface, logger *zap.SugaredLogger, producer sarama.SyncProducer) sender {
	return &core{
		ctx:      ctx,
		user:     user,
		producer: producer,
		logger:   logger,
	}
}

type core struct {
	ctx      context.Context
	user     userPkg.Interface
	producer sarama.SyncProducer
	logger   *zap.SugaredLogger
}

func (c *core) userCreate(msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, brokerDataService)
	defer span.Finish()

	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	c.logger.Debugf("user [%s]", user.String())

	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.StringEncoder(consts.UserCreate),
		Value: sarama.ByteEncoder(fmt.Sprintf("user [%s] created", user.Name)),
	}

	if err := c.user.Create(c.ctx, user); err != nil {
		c.logger.Errorf("user create: %v", err)

		message.Topic = consts.TopicError
		message.Value = sarama.ByteEncoder(fmt.Sprintf("user [%s] creating error: %v", user.Name, err))
		if errInj := helper.InjectSpanIntoMessage(span, message); errInj != nil {
			return errInj
		}

		_, _, errSend := c.producer.SendMessage(message)
		if errSend != nil {
			c.logger.Errorf("user create error send: %v", err)
		}
		return err
	}

	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return err
	}

	part, offset, err := c.producer.SendMessage(message)
	if err != nil {
		c.logger.Errorf("send: %v", err)
		return err
	}
	c.logger.Debugf("part [ %d ] offset [ %d ]", part, offset)

	return nil
}

func (c *core) userUpdate(msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, brokerDataService)
	defer span.Finish()

	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "message unmarshal")
	}

	c.logger.Debugf("user [%s]", user.String())

	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.StringEncoder(consts.UserUpdate),
		Value: sarama.ByteEncoder(fmt.Sprintf("user [%s] updated", user.Name)),
	}

	if err := c.user.Update(c.ctx, user); err != nil {
		c.logger.Errorf("user create: %v", err)

		message.Topic = consts.TopicError
		message.Value = sarama.ByteEncoder(fmt.Sprintf("user [%s] updating error: %v", user.Name, err))
		if errInj := helper.InjectSpanIntoMessage(span, message); errInj != nil {
			return errInj
		}

		_, _, errSend := c.producer.SendMessage(message)
		if errSend != nil {
			c.logger.Errorf("user create error send: %v", err)
		}
		return err
	}

	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return err
	}

	part, offset, err := c.producer.SendMessage(message)
	if err != nil {
		c.logger.Errorf("send: %v", err)
		return err
	}
	c.logger.Debugf("part [ %d ] offset [ %d ]", part, offset)

	return nil
}

func (c *core) userDelete(msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, brokerDataService)
	defer span.Finish()

	name := string(msg.Value)

	c.logger.Debugf("[%s] name: [%s]", name)

	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.StringEncoder(consts.UserDelete),
		Value: sarama.ByteEncoder(fmt.Sprintf("user [%s] removed", name)),
	}

	if err := c.user.Delete(c.ctx, name); err != nil {
		c.logger.Errorf("[%s] user create: %v", err)

		message.Topic = consts.TopicError
		message.Value = sarama.ByteEncoder(fmt.Sprintf("user [%s] removing error: %v", name, err))
		if errInj := helper.InjectSpanIntoMessage(span, message); errInj != nil {
			return errInj
		}

		_, _, errSend := c.producer.SendMessage(message)
		if errSend != nil {
			c.logger.Errorf("user create error send: %v", err)
		}
		return err
	}

	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return err
	}

	part, offset, err := c.producer.SendMessage(message)
	if err != nil {
		c.logger.Errorf("send: %v", err)
		return err
	}
	c.logger.Debugf("part [ %d ] offset [ %d ]", part, offset)

	return nil
}
