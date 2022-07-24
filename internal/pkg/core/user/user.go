package user

import (
	"log"

	"github.com/pkg/errors"

	cachePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/cache"
	localCachePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/cache/local"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
)

var (
	ErrValidation = errors.New("invalid data")
)

type Interface interface {
	Create(user models.User) error
	Update(user models.User) error
	Delete(name string) error
	Get(name string) (models.User, error)
	List() []models.User
}

func MustNew() Interface {
	cache := localCachePkg.New()
	if err := cache.Migrate(); err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	return &core{
		cache: cache,
	}
}

type core struct {
	cache cachePkg.Interface
}

func (c *core) Create(user models.User) error {
	if user.Name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(ErrValidation, "field: [password] cannot be empty")
	}

	return c.cache.Add(user)
}

func (c *core) Update(user models.User) error {
	if user.Name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(ErrValidation, "field: [password] cannot be empty")
	}

	return c.cache.Update(user)
}

func (c *core) Delete(name string) error {
	if name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	return c.cache.Delete(name)
}

func (c *core) Get(name string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	return c.cache.Get(name)
}

func (c *core) List() []models.User {
	return c.cache.List()
}
