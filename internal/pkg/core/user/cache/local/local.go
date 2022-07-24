package local

import (
	"sync"

	"github.com/pkg/errors"

	cachePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/cache"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
)

var (
	ErrUserNotExists = errors.New("user does not exists")
	ErrUserExists    = errors.New("user exists")
)

type cache struct {
	mu   sync.RWMutex
	data map[string]models.User
}

func New() cachePkg.Interface {
	return &cache{
		mu:   sync.RWMutex{},
		data: make(map[string]models.User),
	}
}

func (c *cache) Add(user models.User) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[user.Name]; ok {
		return errors.Wrapf(ErrUserExists, "user-name: [%s]", user.Name)
	}
	c.data[user.Name] = user
	return nil
}

func (c *cache) Update(user models.User) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[user.Name]; !ok {
		return errors.Wrapf(ErrUserNotExists, "user-name: [%s]", user.Name)
	}
	c.data[user.Name] = user
	return nil
}

func (c *cache) Delete(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[name]; !ok {
		return errors.Wrapf(ErrUserNotExists, "user-name: [%s]", name)
	}
	delete(c.data, name)
	return nil
}

func (c *cache) Get(name string) (models.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if user, ok := c.data[name]; !ok {
		return user, errors.Wrapf(ErrUserNotExists, "user-name: [%s]", name)
	} else {
		return user, nil
	}
}

func (c *cache) List() []models.User {
	c.mu.RLock()
	defer c.mu.RUnlock()

	list := make([]models.User, 0, len(c.data))
	for _, user := range c.data {
		list = append(list, user)
	}
	return list
}

func (c *cache) Migrate() error {
	for _, u := range users {
		user := models.User{
			Name:     u.name,
			Password: u.password,
		}
		if err := c.Add(user); err != nil {
			return err
		}
	}
	return nil
}
