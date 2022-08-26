package mailing

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
)

func NewHandler(logger *zap.SugaredLogger, producer sarama.SyncProducer) *Handler {
	return &Handler{
		logger: logger,
		sender: newSender(logger, producer),
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
		switch msg.Topic {
		case consts.TopicMailing:
			session.MarkMessage(msg, "success")
			if err := h.sender.sendSuccess(msg); err != nil {
				return errors.Wrap(err, "send message")
			}
		case consts.TopicError:
			session.MarkMessage(msg, "error")
			if err := h.sender.sendError(msg); err != nil {
				return errors.Wrap(err, "send message")
			}
		}
	}
	return nil
}
