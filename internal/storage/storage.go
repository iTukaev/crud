package storage

import (
	"log"
	"strconv"

	"github.com/pkg/errors"
)

var data map[uint]*User

var UserNotExists = errors.New("user dows not exist")
var UserExists = errors.New("user exists")

func init() {
	log.Println("init storage")
	data = make(map[uint]*User)
	u, _ := NewUser("Kirill", "123456")
	if err := Add(u); err != nil {
		log.Panic(err)
	}
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
		return errors.Wrap(UserExists, strconv.FormatUint(uint64(u.GetId()), 10))
	}
	data[u.GetId()] = u
	return nil
}

func Update(u *User) error {
	if _, ok := data[u.GetId()]; !ok {
		return errors.Wrap(UserNotExists, strconv.FormatUint(uint64(u.GetId()), 10))
	}
	data[u.GetId()] = u
	return nil
}

func Delete(id uint) error {
	if _, ok := data[id]; ok {
		delete(data, id)
		return nil
	}
	return errors.Wrap(UserNotExists, strconv.FormatUint(uint64(id), 10))
}
