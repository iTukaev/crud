package user

import (
	"context"
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
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (models.User, error)
	List(ctx context.Context) ([]models.User, error)
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

func (c *core) Create(ctx context.Context, user models.User) error {
	if user.Name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(ErrValidation, "field: [password] cannot be empty")
	}

	return c.cache.Add(ctx, user)
}

func (c *core) Update(ctx context.Context, user models.User) error {
	if user.Name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(ErrValidation, "field: [password] cannot be empty")
	}

	return c.cache.Update(ctx, user)
}

func (c *core) Delete(ctx context.Context, name string) error {
	if name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	return c.cache.Delete(ctx, name)
}

func (c *core) Get(ctx context.Context, name string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	return c.cache.Get(ctx, name)
}

func (c *core) List(ctx context.Context) ([]models.User, error) {
	return c.cache.List(ctx)
}
