package mailing

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

const (
	mailingService = "mailing"

	expirationCached = 10 * time.Minute
)

type sender interface {
	sendSuccess(ctx context.Context, msg *sarama.ConsumerMessage) error
	sendError(ctx context.Context, msg *sarama.ConsumerMessage) error
}

func newSender(logger *zap.SugaredLogger, producer sarama.SyncProducer, client *redis.Client) sender {
	return &core{
		producer: producer,
		logger:   logger,
		cache:    client,
	}
}

type core struct {
	producer sarama.SyncProducer
	logger   *zap.SugaredLogger
	cache    *redis.Client
}

func (c *core) sendSuccess(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, mailingService)
	defer span.Finish()
	uid, pub := helper.ExtractUidPubFromMessage(msg)

	switch pub {
	case pb.Wait_pub.String():
		if err := c.cache.Publish(ctx, string(msg.Key), msg.Value).Err(); err != nil {
			if err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
				Topic:   consts.TopicMailing,
				Key:     sarama.ByteEncoder(msg.Key),
				Value:   sarama.ByteEncoder(msg.Value),
				Headers: adaptor.ConsumerHeaderToProducer(msg.Headers),
			}); err != nil {
				return err
			}
		}
	case pb.Wait_cache.String():
		if err := c.cache.Set(ctx, uid, msg.Value, expirationCached).Err(); err != nil {
			if err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
				Topic:   consts.TopicMailing,
				Key:     sarama.ByteEncoder(msg.Key),
				Value:   sarama.ByteEncoder(msg.Value),
				Headers: adaptor.ConsumerHeaderToProducer(msg.Headers),
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *core) sendError(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, mailingService)
	defer span.Finish()
	uid, pub := helper.ExtractUidPubFromMessage(msg)

	switch pub {
	case pb.Wait_pub.String():
		if err := c.cache.Set(ctx, string(msg.Key), msg.Value, expirationCached).Err(); err != nil {
			if err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
				Topic:   consts.TopicError,
				Key:     sarama.ByteEncoder(msg.Key),
				Value:   sarama.ByteEncoder(msg.Value),
				Headers: adaptor.ConsumerHeaderToProducer(msg.Headers),
			}); err != nil {
				return err
			}
		}
	case pb.Wait_cache.String():
		if err := c.cache.Publish(ctx, uid, msg.Value).Err(); err != nil {
			if err = c.sendMessageWithCtx(ctx, &sarama.ProducerMessage{
				Topic:   consts.TopicError,
				Key:     sarama.ByteEncoder(msg.Key),
				Value:   sarama.ByteEncoder(msg.Value),
				Headers: adaptor.ConsumerHeaderToProducer(msg.Headers),
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *core) sendMessageWithCtx(ctx context.Context, message *sarama.ProducerMessage) error {
	if err := helper.InjectHeaders(ctx, message); err != nil {
		return err
	}
	_, _, err := c.producer.SendMessage(message)
	return err
}
