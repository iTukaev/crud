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
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

const (
	mailingService = "mailing"
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
	span := helper.GetSpanFromMessage(msg, mailingService)
	defer span.Finish()

	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	if randError() == nil {
		c.logger.Infof("result sended success [%s] [%s]", msg.Key, msg.Value)
		return nil
	}

	message := &sarama.ProducerMessage{
		Topic: consts.TopicMailing,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return err
	}

	_, _, err := c.producer.SendMessage(message)
	if err != nil {
		c.logger.Errorf("send: %v", err)
		return err
	}

	return nil
}

func (c *core) sendError(msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, mailingService)
	defer span.Finish()

	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "message unmarshal")
	}

	if randError() == nil {
		c.logger.Infof("error sended success [%s] [%s]", msg.Key, msg.Value)
		return nil
	}
	message := &sarama.ProducerMessage{
		Topic: consts.TopicError,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return err
	}

	_, _, err := c.producer.SendMessage(message)
	if err != nil {
		c.logger.Errorf("send: %v", err)
		return err
	}

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
