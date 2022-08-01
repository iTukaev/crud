package repo

import (
	"context"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
)

type Interface interface {
	UserCreate(ctx context.Context, user models.User) error
	UserUpdate(ctx context.Context, user models.User) error
	UserDelete(ctx context.Context, name string) error
	UserGet(ctx context.Context, name string) (models.User, error)
	UserList(ctx context.Context, order bool, limit, offset uint32) ([]models.User, error)
	Close()
}
