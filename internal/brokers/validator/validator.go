package validator

import (
	"encoding/json"
	"regexp"

	"github.com/Shopify/sarama"
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
	userCreate(msg *sarama.ConsumerMessage) error
	userUpdate(msg *sarama.ConsumerMessage) error
	userDelete(msg *sarama.ConsumerMessage) error
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

func (c *core) userCreate(msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, validateService)
	defer span.Finish()

	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	c.logger.Debugf("user [%s]", user.String())
	if err := createValidator(user); err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: consts.TopicData,
		Key:   sarama.StringEncoder(consts.UserCreate),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return err
	}

	part, offset, err := c.producer.SendMessage(message)
	if err != nil {
		c.logger.Errorf("send: %v", err)
		return err
	}
	c.logger.Debugf("part [ %d ] offset [ %d ]", part, offset)

	return nil
}

func (c *core) userUpdate(msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, validateService)
	defer span.Finish()

	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "message unmarshal")
	}

	c.logger.Debugf("user [%s]", user.String())
	if err := updateValidator(user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	message := &sarama.ProducerMessage{
		Topic: consts.TopicData,
		Key:   sarama.StringEncoder(consts.UserUpdate),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return err
	}

	part, offset, err := c.producer.SendMessage(message)
	if err != nil {
		c.logger.Errorf("send: %v", err)
		return err
	}
	c.logger.Debugf("part [ %d ] offset [ %d ]", part, offset)

	return nil
}

func (c *core) userDelete(msg *sarama.ConsumerMessage) error {
	span := helper.GetSpanFromMessage(msg, validateService)
	defer span.Finish()

	name := string(msg.Value)

	c.logger.Debugf("name: [%s]", name)
	if err := deleteValidator(name); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	message := &sarama.ProducerMessage{
		Topic: consts.TopicData,
		Key:   sarama.StringEncoder(consts.UserDelete),
		Value: sarama.ByteEncoder(msg.Value),
	}
	if err := helper.InjectSpanIntoMessage(span, message); err != nil {
		return err
	}

	part, offset, err := c.producer.SendMessage(message)
	if err != nil {
		c.logger.Errorf("send: %v", err)
		return err
	}
	c.logger.Debugf("part [ %d ] offset [ %d ]", part, offset)

	return nil
}

func createValidator(user models.User) error {
	if user.Name == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [password] cannot be empty")
	}
	if !email.MatchString(user.Email) {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [email] cannot be empty")
	}
	if user.FullName == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [full_name] cannot be empty")
	}
	return nil
}

func updateValidator(user models.User) error {
	if user.Name == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [password] cannot be empty")
	}
	if !email.MatchString(user.Email) {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [email] cannot be empty")
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
