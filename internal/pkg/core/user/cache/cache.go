package cache

import "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"

type Interface interface {
	Add(user models.User) error
	Update(user models.User) error
	Delete(id string) error
	Get(id string) (models.User, error)
	List() []models.User
	Migrate() error
}
