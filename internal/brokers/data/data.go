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
	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	produceHeaders, meta := headers(msg.Headers)

	c.logger.Debugf("[%s] user [%s]", meta, user.String())

	if err := c.user.Create(c.ctx, user); err != nil {
		c.logger.Errorf("[%s] user create: %v", meta, err)
		_, _, errSend := c.producer.SendMessage(&sarama.ProducerMessage{
			Topic:   consts.TopicError,
			Key:     sarama.StringEncoder(consts.UserCreate),
			Value:   sarama.ByteEncoder(fmt.Sprintf("user [%s] creating error: %v", user.Name, err)),
			Headers: produceHeaders,
		})
		if errSend != nil {
			c.logger.Errorf("[%s] user create error send: %v", meta, err)
		}
		return err
	}

	part, offset, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:   consts.TopicMailing,
		Key:     sarama.StringEncoder(consts.UserCreate),
		Value:   sarama.ByteEncoder(fmt.Sprintf("user [%s] created", user.Name)),
		Headers: produceHeaders,
	})
	if err != nil {
		c.logger.Errorf("[%s] send: %v", meta, err)
		return err
	}
	c.logger.Debugf("[%s] part [ %d ] offset [ %d ]", meta, part, offset)

	return nil
}

func (c *core) userUpdate(msg *sarama.ConsumerMessage) error {
	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "message unmarshal")
	}
	produceHeaders, meta := headers(msg.Headers)

	c.logger.Debugf("[%s] user [%s]", meta, user.String())

	if err := c.user.Update(c.ctx, user); err != nil {
		c.logger.Errorf("[%s] user create: %v", meta, err)
		_, _, errSend := c.producer.SendMessage(&sarama.ProducerMessage{
			Topic:   consts.TopicError,
			Key:     sarama.StringEncoder(consts.UserUpdate),
			Value:   sarama.ByteEncoder(fmt.Sprintf("user [%s] updating error: %v", user.Name, err)),
			Headers: produceHeaders,
		})
		if errSend != nil {
			c.logger.Errorf("[%s] user create error send: %v", meta, err)
		}
		return err
	}

	part, offset, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:   consts.TopicMailing,
		Key:     sarama.StringEncoder(consts.UserUpdate),
		Value:   sarama.ByteEncoder(fmt.Sprintf("user [%s] updated", user.Name)),
		Headers: produceHeaders,
	})
	if err != nil {
		c.logger.Errorf("[%s] send: %v", meta, err)
		return err
	}
	c.logger.Debugf("[%s] part [ %d ] offset [ %d ]", meta, part, offset)

	return nil
}

func (c *core) userDelete(msg *sarama.ConsumerMessage) error {
	name := string(msg.Value)
	produceHeaders, meta := headers(msg.Headers)

	c.logger.Debugf("[%s] name: [%s]", meta, name)
	if err := c.user.Delete(c.ctx, name); err != nil {
		c.logger.Errorf("[%s] user create: %v", meta, err)
		_, _, errSend := c.producer.SendMessage(&sarama.ProducerMessage{
			Topic:   consts.TopicError,
			Key:     sarama.StringEncoder(consts.UserDelete),
			Value:   sarama.ByteEncoder(fmt.Sprintf("user [%s] removing error: %v", name, err)),
			Headers: produceHeaders,
		})
		if errSend != nil {
			c.logger.Errorf("[%s] user create error send: %v", meta, err)
		}
		return err
	}

	part, offset, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:   consts.TopicMailing,
		Key:     sarama.StringEncoder(consts.UserDelete),
		Value:   sarama.ByteEncoder(fmt.Sprintf("user [%s] deleted", name)),
		Headers: produceHeaders,
	})
	if err != nil {
		c.logger.Errorf("[%s] send: %v", meta, err)
		return err
	}
	c.logger.Debugf("[%s] part [ %d ] offset [ %d ]", meta, part, offset)

	return nil
}

func headers(headers []*sarama.RecordHeader) ([]sarama.RecordHeader, string) {
	res := make([]sarama.RecordHeader, 0, len(headers))
	var meta string
	for _, h := range headers {
		if string(h.Key) == "meta" {
			meta = string(h.Value)
		}
		res = append(res, *h)
	}
	return res, meta
}
