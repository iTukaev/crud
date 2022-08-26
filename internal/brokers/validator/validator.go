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
	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	produceHeaders, meta := headers(msg.Headers)

	c.logger.Debugf("[%s] user [%s]", meta, user.String())
	if err := createValidator(user); err != nil {
		return err
	}

	part, offset, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:   consts.TopicData,
		Key:     sarama.StringEncoder(consts.UserCreate),
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

func (c *core) userUpdate(msg *sarama.ConsumerMessage) error {
	var user models.User
	if err := json.Unmarshal(msg.Value, &user); err != nil {
		return errors.Wrap(err, "message unmarshal")
	}
	produceHeaders, meta := headers(msg.Headers)

	c.logger.Debugf("[%s] user [%s]", meta, user.String())
	if err := updateValidator(user); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	part, offset, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:   consts.TopicData,
		Key:     sarama.StringEncoder(consts.UserUpdate),
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

func (c *core) userDelete(msg *sarama.ConsumerMessage) error {
	name := string(msg.Value)
	produceHeaders, meta := headers(msg.Headers)

	c.logger.Debugf("[%s] name: [%s]", meta, name)
	if err := deleteValidator(name); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	part, offset, err := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic:   consts.TopicData,
		Key:     sarama.StringEncoder(consts.UserDelete),
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
	return nil
}

func deleteValidator(name string) error {
	if name == "" {
		return errors.Wrap(errorsPkg.ErrValidation, "field: [name] cannot be empty")
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
