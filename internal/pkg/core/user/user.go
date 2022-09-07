//go:generate mockgen -source=user.go -destination=./mock/user_mock.go -package=mock

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/counter"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/customerrors"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
)

const (
	ctxTimeout     = 5 * time.Second
	expirationTime = 1 * time.Minute
)

type Interface interface {
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (models.User, error)
	List(ctx context.Context, order bool, limit, offset uint64) ([]models.User, error)
	Data(ctx context.Context, uid string) ([]byte, error)
}

func New(data repoPkg.Interface, logger *zap.SugaredLogger, client *redis.Client) Interface {
	return &core{
		data:   data,
		logger: logger,
		cache:  client,
	}
}

type core struct {
	data   repoPkg.Interface
	logger *zap.SugaredLogger
	cache  *redis.Client
}

func (c *core) Create(ctx context.Context, user models.User) error {
	c.logger.Debugln("Create", user)
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if _, err := c.cache.Get(ctx, user.Name).Bytes(); err == nil {
		counter.Hit.Inc()
		return errorsPkg.ErrUserAlreadyExists
	}

	if _, err := c.data.UserGet(ctx, user.Name); err == nil {
		counter.Miss.Inc()
		return errorsPkg.ErrUserAlreadyExists
	} else if !errors.Is(err, errorsPkg.ErrUserNotFound) {
		return err
	}
	if err := c.data.UserCreate(ctx, user); err != nil {
		return err
	}

	return nil
}

func (c *core) Update(ctx context.Context, user models.User) error {
	c.logger.Debugln("Update", user)
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	old, err := c.data.UserGet(ctx, user.Name)
	if err != nil {
		return err
	}
	if err = c.data.UserUpdate(ctx, user); err != nil {
		return err
	}

	user.CreatedAt = old.CreatedAt
	if err = c.cache.Set(ctx, user.Name, &user, expirationTime).Err(); err != nil {
		c.logger.Errorf("set to cache: %v", err)
	}

	return nil
}

func (c *core) Delete(ctx context.Context, name string) error {
	c.logger.Debugln("Delete", name)
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if _, err := c.data.UserGet(ctx, name); err != nil {
		return err
	}
	if err := c.data.UserDelete(ctx, name); err != nil {
		return err
	}

	if err := c.cache.Del(ctx, name).Err(); err != nil {
		if !errors.Is(err, redis.Nil) {
			c.logger.Errorf("remove from cache: %v", err)
		}
	}

	return nil
}

func (c *core) Get(ctx context.Context, name string) (models.User, error) {
	c.logger.Debugln("Get", name)
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if data, err := c.cache.Get(ctx, name).Bytes(); err == nil {
		counter.Hit.Inc()
		var user models.User
		if err = json.Unmarshal(data, &user); err == nil {
			return user, nil
		}
		c.logger.Errorf("unmarshal cached data: %v", err)
	}

	counter.Miss.Inc()
	user, err := c.data.UserGet(ctx, name)
	if err != nil {
		return user, err
	}
	if err = c.cache.Set(ctx, name, &user, expirationTime).Err(); err != nil {
		c.logger.Errorf("set user to cache: %v", err)
	}

	return user, nil
}

func (c *core) List(ctx context.Context, order bool, limit, offset uint64) ([]models.User, error) {
	c.logger.Debugln("List", order, limit, offset)
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	key := fmt.Sprintf("%v_%d_%d", order, limit, offset)
	if data, err := c.cache.Get(ctx, key).Bytes(); err == nil {
		counter.Hit.Inc()
		users := make([]models.User, 0)
		if err = json.Unmarshal(data, &users); err == nil {
			return users, nil
		}
		c.logger.Errorf("unmarshal cached data: %v", err)
	}

	counter.Miss.Inc()
	users, err := c.data.UserList(ctx, order, limit, offset)
	if err != nil {
		return users, err
	}
	if err = c.cache.Set(ctx, key, users, expirationTime).Err(); err != nil {
		c.logger.Errorf("set users list to cache: %v", err)
	}

	return users, nil
}

func (c *core) Data(ctx context.Context, uid string) ([]byte, error) {
	c.logger.Debugln("Data", uid)

	return c.cache.Get(ctx, uid).Bytes()
}
