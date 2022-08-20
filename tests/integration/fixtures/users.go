//go:build integration
// +build integration

package fixtures

import "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"

var (
	User1 = models.NewUser().
		NameSet("Ivan").
		PasswordSet("123").
		EmailSet("ivan@email.com").
		FullNameSet("Ivan the Dummy").
		CreatedAtSet(1234567890)

	User2 = models.NewUser().
		NameSet("Miron").
		PasswordSet("123").
		EmailSet("miron@email.com").
		FullNameSet("Miron the Simple in the field").
		CreatedAtSet(1234567890)

	ExistedUser2 = models.NewUser().
			NameSet("Piter").
			PasswordSet("123").
			EmailSet("piter@email.com").
			FullNameSet("Piter Parker").
			CreatedAtSet(1659447420)

	ExistedUser1 = models.NewUser().
			NameSet("Berta").
			PasswordSet("654").
			EmailSet("berta@email.com").
			FullNameSet("Big Berta").
			CreatedAtSet(1659447450)
)
