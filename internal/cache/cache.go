package cache

import (
	"context"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
)

type Interface interface {
	Set(ctx context.Context, user models.User)
	Get(ctx context.Context, name string) (models.User, bool)
	Delete(ctx context.Context, name string)
}
