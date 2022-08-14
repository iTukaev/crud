//go:generate mockgen -source=user.go -destination=./mock/user_mock.go -package=mock

package user

import (
	"context"

	"github.com/pkg/errors"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
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
	if _, err := c.data.UserGet(ctx, user.Name); err == nil {
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
	if _, err := c.data.UserGet(ctx, user.Name); err != nil {
		return err
	}
	if err := c.data.UserUpdate(ctx, user); err != nil {
		return err
	}

	return nil
}

func (c *core) Delete(ctx context.Context, name string) error {
	if _, err := c.data.UserGet(ctx, name); err != nil {
		return err
	}
	if err := c.data.UserDelete(ctx, name); err != nil {
		return err
	}

	return nil
}

func (c *core) Get(ctx context.Context, name string) (models.User, error) {
	return c.data.UserGet(ctx, name)
}

func (c *core) List(ctx context.Context, order bool, limit, offset uint64) ([]models.User, error) {
	return c.data.UserList(ctx, order, limit, offset)
}
