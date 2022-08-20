//go:generate chaingen models.go

package models

import (
	"fmt"
	"time"
)

type User struct {
	Name      string `db:"name"`
	Password  string `db:"password"`
	Email     string `db:"email"`
	FullName  string `db:"full_name"`
	CreatedAt int64  `db:"created_at"`
}

func (u User) String() string {
	return fmt.Sprintf("name: [%s], full_name: [%s], email: [%s], created_at: [%v]",
		u.Name, u.FullName, u.Email, time.Unix(u.CreatedAt, 0))
}
