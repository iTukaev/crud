package user

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
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

func MustNew(data repoPkg.Interface) Interface {
	return &core{
		data: data,
	}
}

type core struct {
	data repoPkg.Interface
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

	if _, err := c.data.UserGet(ctx, user.Name); err == nil {
		return errorsPkg.ErrUserAlreadyExists
	}
	if err := c.data.UserCreate(ctx, user); err != nil {
		return err
	}

	return nil
}

func (c *core) Update(ctx context.Context, user models.User) error {
	if user.Name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}

	if _, err := c.data.UserGet(ctx, user.Name); err != nil {
		return err
	}
	if err := c.data.UserUpdate(ctx, user); err != nil {
		return err
	}

	return nil
}

func (c *core) Delete(ctx context.Context, name string) error {
	if name == "" {
		return errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}

	if _, err := c.data.UserGet(ctx, name); err != nil {
		return err
	}
	if err := c.data.UserDelete(ctx, name); err != nil {
		return err
	}

	return nil
}

func (c *core) Get(ctx context.Context, name string) (models.User, error) {
	if name == "" {
		return models.User{}, errors.Wrap(ErrValidation, "field: [name] cannot be empty")
	}

	return c.data.UserGet(ctx, name)
}

func (c *core) List(ctx context.Context, order bool, limit, offset uint64) ([]models.User, error) {
	return c.data.UserList(ctx, order, limit, offset)
}
