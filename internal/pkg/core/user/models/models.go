//go:generate chaingen models.go

package models

import (
	"fmt"
	"time"
)

type User struct {
	Name      string `json:"name" db:"name"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	FullName  string `json:"full_name" db:"full_name"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}

func (u *User) String() string {
	return fmt.Sprintf("name: [%s], full_name: [%s], email: [%s], created_at: [%v]",
		u.Name, u.FullName, u.Email, time.Unix(u.CreatedAt, 0))
}

type UserListParams struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
	Order  bool   `json:"order"`
}
