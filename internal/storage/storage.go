package storage

import (
	"github.com/pkg/errors"
	"log"
)

var data map[uint]*User

var UserNotExists = errors.New("user does not exists")
var UserExists = errors.New("user exists")

func Init() error {
	log.Println("init storage")
	data = make(map[uint]*User)
	if err := migrateUsers(); err != nil {
		return errors.Wrap(err, "migration failed")
	}
	return nil
}

func List() []*User {
	res := make([]*User, 0, len(data))
	for _, v := range data {
		res = append(res, v)
	}
	return res
}

func Add(u *User) error {
	if _, ok := data[u.GetId()]; ok {
		return errors.Wrapf(UserExists, "%d", u.GetId())
	}
	data[u.GetId()] = u
	return nil
}

func Update(id uint, name, pwd string) (*User, error) {
	if _, ok := data[id]; !ok {
		return nil, errors.Wrapf(UserNotExists, "%d", id)
	}
	u := &User{}
	_ = u.SetId(id)
	if err := u.SetName(name); err != nil {
		return nil, err
	}
	if err := u.SetPwd(pwd); err != nil {
		return nil, err
	}

	data[id] = u
	return u, nil
}

func Delete(id uint) (*User, error) {
	if u, ok := data[id]; !ok {
		return nil, errors.Wrapf(UserNotExists, "%d", id)
	} else {
		delete(data, id)
		return u, nil
	}
}
