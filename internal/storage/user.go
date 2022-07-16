package storage

import (
	"fmt"
)

var lastId = uint(0)

type User struct {
	id       uint
	name     string
	password string
}

func NewUser(name, passwprd string) (*User, error) {
	u := User{}
	if err := u.SetName(name); err != nil {
		return nil, err
	}
	if err := u.SetPassword(passwprd); err != nil {
		return nil, err
	}
	lastId++
	u.id = lastId
	return &u, nil
}

func (u *User) SetName(name string) error {
	if len(name) == 0 || len(name) > 10 {
		return fmt.Errorf("bad name <%v>", name)
	}
	u.name = name
	return nil
}

func (u *User) SetPassword(pwd string) error {
	if len(pwd) < 6 || len(pwd) > 10 {
		return fmt.Errorf("bad password <%v>", pwd)
	}
	u.password = pwd
	return nil
}

func (u User) String() string {
	return fmt.Sprintf("%d: %s / %s", u.id, u.name, u.password)
}

func (u User) GetName() string {
	return u.name
}

func (u User) GetPassword() string {
	return u.password
}

func (u User) GetId() uint {
	return u.id
}
