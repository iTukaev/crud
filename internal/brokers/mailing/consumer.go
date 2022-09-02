package mailing

import (
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

func NewHandler(logger *zap.SugaredLogger, producer sarama.SyncProducer, client *redis.Client) *Handler {
	return &Handler{
		logger: logger,
		sender: newSender(logger, producer, client),
	}
}

type Handler struct {
	logger *zap.SugaredLogger
	sender sender
}

func (h *Handler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		return h.handleMessage(session, msg)
	}
	return nil
}

func (h *Handler) handleMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	uid, pub := helper.ExtractUidPubFromMessage(msg)
	ctx := helper.InjectUidPubToCtx(session.Context(), uid, pub)

	switch msg.Topic {
	case consts.TopicMailing:
		session.MarkMessage(msg, "success")
		if err := h.sender.sendSuccess(ctx, msg); err != nil {
			return errors.Wrap(err, "send message")
		}
	case consts.TopicError:
		session.MarkMessage(msg, "error")
		if err := h.sender.sendError(ctx, msg); err != nil {
			return errors.Wrap(err, "send message")
		}
	}
	return nil
}
