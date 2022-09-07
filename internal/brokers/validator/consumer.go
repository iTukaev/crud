package validator

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/customerrors"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
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
		return h.handleMessage(session, msg)
	}
	return nil
}

func (h *Handler) handleMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) error {
	uid, pub := helper.ExtractUidPubFromMessage(msg)
	ctx := helper.InjectUidPubToCtx(session.Context(), uid, pub)

	switch string(msg.Key) {
	case consts.UserCreate:
		if err := h.sender.userCreate(ctx, msg); err != nil {
			return errors.Wrap(err, "user create")
		}
		session.MarkMessage(msg, "")
	case consts.UserUpdate:
		if err := h.sender.userUpdate(ctx, msg); err != nil {
			return errors.Wrap(err, "user update")
		}
		session.MarkMessage(msg, "")
	case consts.UserDelete:
		if err := h.sender.userDelete(ctx, msg); err != nil {
			return errors.Wrap(err, "user delete")
		}
		session.MarkMessage(msg, "")
	case consts.UserGet:
		if err := h.sender.userGet(ctx, msg); err != nil {
			return errors.Wrap(err, "user get")
		}
		session.MarkMessage(msg, "")
	default:
		session.MarkMessage(msg, "invalid_key")
		return errors.Wrap(errorsPkg.ErrValidation, "invalid message key")
	}
	return nil
}
