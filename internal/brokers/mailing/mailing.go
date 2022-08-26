package mailing

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
)

type sender interface {
	sendSuccess(msg *sarama.ConsumerMessage) error
	sendError(msg *sarama.ConsumerMessage) error
}

func newSender(logger *zap.SugaredLogger, producer sarama.SyncProducer) sender {
	return &core{
		producer: producer,
		logger:   logger,
	}
}

type core struct {
	producer sarama.SyncProducer
	logger   *zap.SugaredLogger
}

func (c *core) sendSuccess(msg *sarama.ConsumerMessage) error {
	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	produceHeaders, meta := headers(msg.Headers)

	if randError() == nil {
		c.logger.Infof("[%s] result sended success [%s] [%s]", meta, msg.Key, msg.Value)
		return nil
	}
	_, _, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:   consts.TopicMailing,
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Headers: produceHeaders,
	})
	if err != nil {
		c.logger.Errorf("[%s] send: %v", meta, err)
		return err
	}

	return nil
}

func (c *core) sendError(msg *sarama.ConsumerMessage) error {
	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "message unmarshal")
	}
	produceHeaders, meta := headers(msg.Headers)

	if randError() == nil {
		c.logger.Infof("[%s] error sended success [%s] [%s]", meta, msg.Key, msg.Value)
		return nil
	}
	part, offset, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:   consts.TopicError,
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Headers: produceHeaders,
	})
	if err != nil {
		c.logger.Errorf("[%s] send: %v", meta, err)
		return err
	}
	c.logger.Debugf("[%s] part [ %d ] offset [ %d ]", meta, part, offset)

	return nil
}

func randError() error {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	if rand.Intn(20) > 15 {
		return errors.New("sending error")
	}
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
