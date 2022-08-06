package user

import (
	"context"
	"time"

	"github.com/pkg/errors"

	cachePkg "gitlab.ozon.dev/iTukaev/homework/internal/cache"
	localCachePkg "gitlab.ozon.dev/iTukaev/homework/internal/cache/local"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
	postgresPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres"
	pgModels "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres/models"
)

var (
	ErrValidation = errors.New("invalid data")
)

type Interface interface {
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (models.User, error)
	List(ctx context.Context, order bool, limit, offset uint64) ([]models.User, error)
}

func MustNew(ctx context.Context, pg pgModels.Config) Interface {
	return &core{
		db:    postgresPkg.MustNew(ctx, pg.Host, pg.Port, pg.User, pg.Password, pg.DBName),
		cache: localCachePkg.New(),
	}
}

type core struct {
	db    repoPkg.Interface
	cache cachePkg.Interface
}

func (c *core) Create(ctx context.Context, user models.User) error {
	if user.Name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}
	if user.Password == "" {
		return errors.Wrap(ErrValidation, "field: [password] cannot be empty")
	}
	if user.Email == "" {
		return errors.Wrap(ErrValidation, "field: [email] cannot be empty")
	}
	if user.FullName == "" {
		return errors.Wrap(ErrValidation, "field: [full_name] cannot be empty")
	}
	user.CreatedAt = time.Now().Unix()

	if u, _ := c.db.UserGet(ctx, user.Name); u.Name == user.Name {
		return localCachePkg.ErrUserAlreadyExists
	}
	if err := c.db.UserCreate(ctx, user); err != nil {
		return err
	}

	c.cache.Set(ctx, user)
	return nil
}

func (c *core) Update(ctx context.Context, user models.User) error {
	if user.Name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}

	if u, _ := c.db.UserGet(ctx, user.Name); u.Name != user.Name {
		return localCachePkg.ErrUserNotFound
	}
	if err := c.db.UserUpdate(ctx, user); err != nil {
		return err
	}

	c.cache.Set(ctx, user)
	return nil
}

func (c *core) Delete(ctx context.Context, name string) error {
	if name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}

	if err := c.db.UserDelete(ctx, name); err != nil {
		return err
	}

	c.cache.Delete(ctx, name)
	return nil
}

func (c *core) Get(ctx context.Context, name string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}

	if user, ok := c.cache.Get(ctx, name); ok {
		return user, nil
	}

	return c.db.UserGet(ctx, name)
}

func (c *core) List(ctx context.Context, order bool, limit, offset uint64) ([]models.User, error) {
	return c.db.UserList(ctx, order, limit, offset)
}
