package local

import (
	"context"
	"log"
	"sort"
	"sync"

	"github.com/pkg/errors"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

func New(workersCount int, logger loggerPkg.Interface) repoPkg.Interface {
	log.Println("With local storage started")
	return &cache{
		mu:     sync.RWMutex{},
		data:   make(map[string]models.User),
		poolCh: make(chan struct{}, workersCount),
		logger: logger,
	}
}

type cache struct {
	mu     sync.RWMutex
	data   map[string]models.User
	poolCh chan struct{}
	logger loggerPkg.Interface
}

func (c *cache) UserCreate(ctx context.Context, user models.User) error {
	c.logger.Debug("UserCreate, cached func", user.String())
	select {
	case <-ctx.Done():
		return errorsPkg.ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		c.data[user.Name] = user
		return nil
	}
}

func (c *cache) UserUpdate(ctx context.Context, user models.User) error {
	c.logger.Debug("UserUpdate, cached func", user.String())
	select {
	case <-ctx.Done():
		return errorsPkg.ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		u := c.data[user.Name]
		if user.Email != "" {
			u.Email = user.Email
		}
		if user.Password != "" {
			u.Password = user.Password
		}
		if user.FullName != "" {
			u.FullName = user.FullName
		}

		c.data[user.Name] = u
		return nil
	}
}

func (c *cache) UserDelete(ctx context.Context, name string) error {
	c.logger.Debug("UserDelete, cached func", name)
	select {
	case <-ctx.Done():
		return errorsPkg.ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.Lock()
		defer func() {
			c.mu.Unlock()
			<-c.poolCh
		}()

		delete(c.data, name)
		return nil
	}
}

func (c *cache) UserGet(ctx context.Context, name string) (models.User, error) {
	c.logger.Debug("UserGet, cached func", name)
	select {
	case <-ctx.Done():
		return models.User{}, errorsPkg.ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.RLock()
		defer func() {
			c.mu.RUnlock()
			<-c.poolCh
		}()

		if user, ok := c.data[name]; !ok {
			return user, errors.Wrapf(errorsPkg.ErrUserNotFound, "user-name: [%s]", name)
		} else {
			return user, nil
		}
	}
}

func (c *cache) UserList(ctx context.Context, order bool, limit, offset uint64) ([]models.User, error) {
	c.logger.Debug("UserList, cached func", order, limit, offset)
	select {
	case <-ctx.Done():
		return nil, errorsPkg.ErrTimeout
	case c.poolCh <- struct{}{}:
		c.mu.RLock()
		defer func() {
			c.mu.RUnlock()
			<-c.poolCh
		}()

		if len(c.data) < int(limit*offset) {
			return make([]models.User, 0), nil
		}

		list := make([]models.User, 0, len(c.data))
		for _, user := range c.data {
			list = append(list, user)
		}

		sort.Slice(list, func(i, j int) bool {
			if order {
				return list[i].Name > list[j].Name
			}
			return list[i].Name < list[j].Name
		})

		min := limit * offset
		if len(list) < int(limit*(offset+1)) {
			return list[min:], nil
		} else {
			max := limit * (offset + 1)
			return list[min:max], nil
		}
	}
}

func (c *cache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = nil
	close(c.poolCh)
	c.logger.Info("Cache cleaned")
}
