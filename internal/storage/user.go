package storage

import "fmt"

var lastId = uint(0)

type User struct {
	id       uint
	name     string
	password string
}

func NewUser(name, password string) (*User, error) {
	u := User{}
	if err := u.SetName(name); err != nil {
		return nil, err
	}
	if err := u.SetPwd(password); err != nil {
		return nil, err
	}
	u.id = lastId
	lastId++
	return &u, nil
}

func (u *User) SetId(id uint) error {
	if id < 0 {
		return fmt.Errorf("bad id <%d>", id)
	}
	u.id = id
	return nil
}

func (u *User) SetName(name string) error {
	if len(name) == 0 || len(name) > 10 {
		return fmt.Errorf("bad name <%s>", name)
	}
	u.name = name
	return nil
}

func (u *User) SetPwd(pwd string) error {
	if len(pwd) == 0 || len(pwd) > 10 {
		return fmt.Errorf("bad password <%s>", pwd)
	}
	u.password = pwd
	return nil
}

func (u *User) String() string {
	return fmt.Sprintf("%d: %s / %s", u.id, u.name, u.password)
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetPwd() string {
	return u.password
}

func (u *User) GetId() uint {
	return u.id
}
