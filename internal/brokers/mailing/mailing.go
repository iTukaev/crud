package mailing

import (
	"context"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

const (
	mailingService = "mailing"
)

type sender interface {
	sendSuccess(ctx context.Context, msg *sarama.ConsumerMessage) error
	sendError(ctx context.Context, msg *sarama.ConsumerMessage) error
}

func newSender(logger *zap.SugaredLogger, producer sarama.SyncProducer) sender {
	return &core{
		producer: producer,
		logger:   logger,
	}
}

type core struct {
	// todo: implement Redis
	producer sarama.SyncProducer
	logger   *zap.SugaredLogger
}

func (c *core) sendSuccess(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, mailingService)
	defer span.Finish()
	uid, pub := helper.ExtractUidPubFromMessage(msg)
	_ = uid
	_ = pub
	switch string(msg.Key) {
	case consts.UserGet:
		//todo: implement
	case consts.UserList:
		//todo: implement
	default:
		//todo: implement
	}
	//todo: sending to Redis logic

	message := &sarama.ProducerMessage{
		Topic:   consts.TopicError,
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Headers: adaptor.ConsumerHeaderToProducer(msg.Headers),
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) sendError(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, mailingService)
	defer span.Finish()
	uid, pub := helper.ExtractUidPubFromMessage(msg)
	_ = uid
	_ = pub

	//todo: sending to Redis logic

	message := &sarama.ProducerMessage{
		Topic:   consts.TopicError,
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Headers: adaptor.ConsumerHeaderToProducer(msg.Headers),
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) sendMessageWithCtx(ctx context.Context, message *sarama.ProducerMessage) error {
	if err := helper.InjectHeaders(ctx, message); err != nil {
		return err
	}
	_, _, err := c.producer.SendMessage(message)
	return err
}
