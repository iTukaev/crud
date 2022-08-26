package data

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/customerrors"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
)

func NewHandler(ctx context.Context, user userPkg.Interface, logger *zap.SugaredLogger, producer sarama.SyncProducer) *Handler {
	return &Handler{
		logger: logger,
		sender: newSender(ctx, user, logger, producer),
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
		key := string(msg.Key)
		switch key {
		case consts.UserCreate:
			session.MarkMessage(msg, "create")
			if err := h.sender.userCreate(msg); err != nil {
				return errors.Wrap(err, "user create")
			}
		case consts.UserUpdate:
			session.MarkMessage(msg, "update")
			if err := h.sender.userUpdate(msg); err != nil {
				return errors.Wrap(err, "user update")
			}
		case consts.UserDelete:
			session.MarkMessage(msg, "delete")
			if err := h.sender.userUpdate(msg); err != nil {
				return errors.Wrap(err, "user delete")
			}
		default:
			session.MarkMessage(msg, "invalid_key")
			return errors.Wrap(errorsPkg.ErrValidation, "invalid message key")
		}
	}
	return nil
}
