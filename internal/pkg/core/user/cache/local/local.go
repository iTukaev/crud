package local

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	cachePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/cache"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
)

const (
	workersCount = 10

	contextTimeout = 5 * time.Second
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTimeout           = errors.New("deadline exceeded")
)

func New() cachePkg.Interface {
	return &cache{
		mu:     sync.RWMutex{},
		data:   make(map[string]models.User),
		poolCh: make(chan struct{}, workersCount),
	}
}

type cache struct {
	mu     sync.RWMutex
	data   map[string]models.User
	poolCh chan struct{}
}

func (c *cache) Add(ctx context.Context, user models.User) error {
	select {
	case <-ctx.Done():
		return ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		if _, ok := c.data[user.Name]; ok {
			return errors.Wrapf(ErrUserAlreadyExists, "user-name: [%s]", user.Name)
		}
		c.data[user.Name] = user
		return nil
	}
}

func (c *cache) Update(ctx context.Context, user models.User) error {
	select {
	case <-ctx.Done():
		return ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		if _, ok := c.data[user.Name]; !ok {
			return errors.Wrapf(ErrUserNotFound, "user-name: [%s]", user.Name)
		}
		c.data[user.Name] = user
		return nil
	}
}

func (c *cache) Delete(ctx context.Context, name string) error {
	select {
	case <-ctx.Done():
		return ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		if _, ok := c.data[name]; !ok {
			return errors.Wrapf(ErrUserNotFound, "user-name: [%s]", name)
		}
		delete(c.data, name)
		return nil
	}
}

func (c *cache) Get(ctx context.Context, name string) (models.User, error) {
	select {
	case <-ctx.Done():
		return models.User{}, ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.RLock()
		defer func() {
			c.mu.RUnlock()
			<-c.poolCh
		}()

		if user, ok := c.data[name]; !ok {
			return user, errors.Wrapf(ErrUserNotFound, "user-name: [%s]", name)
		} else {
			return user, nil
		}
	}
}

func (c *cache) List(ctx context.Context) ([]models.User, error) {
	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.RLock()
		defer func() {
			c.mu.RUnlock()
			<-c.poolCh
		}()

		list := make([]models.User, 0, len(c.data))
		for _, user := range c.data {
			list = append(list, user)
		}
		return list, nil
	}
}

func (c *cache) Migrate() error {
	for _, u := range users {
		user := models.User{
			Name:     u.name,
			Password: u.password,
		}
		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		if err := c.Add(ctx, user); err != nil {
			cancel()
			return err
		}
		cancel()
	}
	return nil
}
