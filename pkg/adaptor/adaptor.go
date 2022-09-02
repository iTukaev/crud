package adaptor

import (
	coreModels "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	pbModels "gitlab.ozon.dev/iTukaev/homework/pkg/api/models"
)

func ToUserPbModel(u coreModels.User) *pbModels.User {
	return &pbModels.User{
		Name:      u.Name,
		Password:  u.Password,
		Email:     u.Email,
		FullName:  u.FullName,
		CreatedAt: u.CreatedAt,
	}
}

func ToUserCoreModel(u *pbModels.User) *coreModels.User {
	return &coreModels.User{
		Name:      u.Name,
		Password:  u.Password,
		Email:     u.Email,
		FullName:  u.FullName,
		CreatedAt: u.CreatedAt,
	}
}

func ToUserListPbModel(users []coreModels.User) []*pbModels.User {
	list := make([]*pbModels.User, 0, len(users))
	for _, user := range users {
		list = append(list, ToUserPbModel(user))
	}

	return list
}
