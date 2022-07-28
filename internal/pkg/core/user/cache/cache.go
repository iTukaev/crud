package cache

import (
	"context"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
)

type Interface interface {
	Add(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Migrate() error
}
