package local

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	cachePkg "gitlab.ozon.dev/iTukaev/homework/internal/cache"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
)

const (
	workersCount = 10
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

func (c *cache) Set(ctx context.Context, user models.User) {
	select {
	case <-ctx.Done():
		return
	case c.poolCh <- struct{}{}:
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		c.data[user.Name] = user
		return
	}
}

func (c *cache) Delete(ctx context.Context, name string) {
	select {
	case <-ctx.Done():
		return
	case c.poolCh <- struct{}{}:
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		delete(c.data, name)
		return
	}
}

func (c *cache) Get(ctx context.Context, name string) (models.User, bool) {
	select {
	case <-ctx.Done():
		return models.User{}, false
	case c.poolCh <- struct{}{}:
		c.mu.RLock()
		defer func() {
			c.mu.RUnlock()
			<-c.poolCh
		}()

		user, ok := c.data[name]
		return user, ok
	}
}
