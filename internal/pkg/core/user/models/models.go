package models

import "fmt"

type User struct {
	Name     string
	Password string
}

func (u User) String() string {
	return fmt.Sprintf("name: [%s], password: [%s]", u.Name, u.Password)
}
