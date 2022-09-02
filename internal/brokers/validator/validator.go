package validator

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/customerrors"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	"gitlab.ozon.dev/iTukaev/homework/pkg/helper"
)

const (
	validateService = "validate"
)

var (
	email = regexp.MustCompile(`^.+@[A-Za-z0-9\-_\.]+$`)
)

type sender interface {
	userCreate(ctx context.Context, msg *sarama.ConsumerMessage) error
	userUpdate(ctx context.Context, msg *sarama.ConsumerMessage) error
	userDelete(ctx context.Context, msg *sarama.ConsumerMessage) error
	userGet(ctx context.Context, msg *sarama.ConsumerMessage) error
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

func (c *core) userCreate(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, validateService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	user := models.NewUser()
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	c.logger.Debugf("user [%s]", user.String())

	message := &sarama.ProducerMessage{
		Topic: consts.TopicData,
		Key:   sarama.StringEncoder(consts.UserCreate),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := createValidator(user); err != nil {
		return c.sendValidationErrorWithCtx(ctx, message, err.Error())
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) userUpdate(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, validateService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	user := models.NewUser()
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "message unmarshal")
	}

	c.logger.Debugf("user [%s]", user.String())

	message := &sarama.ProducerMessage{
		Topic: consts.TopicData,
		Key:   sarama.StringEncoder(consts.UserUpdate),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := updateValidator(user); err != nil {
		return c.sendValidationErrorWithCtx(ctx, message, err.Error())
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) userDelete(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, validateService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	name := string(msg.Value)

	message := &sarama.ProducerMessage{
		Topic: consts.TopicData,
		Key:   sarama.StringEncoder(consts.UserDelete),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := deleteValidator(name); err != nil {
		return c.sendValidationErrorWithCtx(ctx, message, err.Error())
	}

	return c.sendMessageWithCtx(ctx, message)
}

func (c *core) userGet(ctx context.Context, msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, validateService)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	name := string(msg.Value)

	message := &sarama.ProducerMessage{
		Topic: consts.TopicData,
		Key:   sarama.StringEncoder(consts.UserGet),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := getValidator(name); err != nil {
		return c.sendValidationErrorWithCtx(ctx, message, err.Error())
	}

	return c.sendMessageWithCtx(ctx, message)
}

func createValidator(user *models.User) error {
	if user.Name == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [password] cannot be empty")
	}
	if !email.MatchString(user.Email) {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [email] has invalid format")
	}
	if user.FullName == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [full_name] cannot be empty")
	}
	return nil
}

func updateValidator(user *models.User) error {
	if user.Name == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [password] cannot be empty")
	}
	if !email.MatchString(user.Email) {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [email] has invalid format")
	}
	if user.FullName == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [full_name] cannot be empty")
	}
	return nil
}

func deleteValidator(name string) error {
	if name == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [name] cannot be empty")
	}
	return nil
}

func getValidator(name string) error {
	if name == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [name] cannot be empty")
	}
	return nil
}

func (c *core) sendValidationErrorWithCtx(
	ctx context.Context,
	message *sarama.ProducerMessage,
	description string,
) error {
	if err := helper.InjectHeaders(ctx, message); err != nil {
		return err
	}
	message.Topic = consts.TopicError
	message.Value = sarama.StringEncoder(description)

	_, _, err := c.producer.SendMessage(message)
	return err
}

func (c *core) sendMessageWithCtx(ctx context.Context, message *sarama.ProducerMessage) error {
	if err := helper.InjectHeaders(ctx, message); err != nil {
		return err
	}
	_, _, err := c.producer.SendMessage(message)
	return err
}
