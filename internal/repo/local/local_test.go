package local

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/customerrors"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

var (
	user1 = models.User{
		Name:      "Ivan",
		Password:  "123",
		Email:     "ivan@email.com",
		FullName:  "Ivan the Dummy",
		CreatedAt: 1660412940,
	}
	user2 = models.User{
		Name:      "Ivan",
		Password:  "123456",
		Email:     "ivanivan@email.com",
		FullName:  "Ivan the Smart guy",
		CreatedAt: 1660412940,
	}
	user3 = models.User{
		Name:      "Boris",
		Password:  "321",
		Email:     "boris@email.com",
		FullName:  "Boris The Blade",
		CreatedAt: 1660412960,
	}
	user4 = models.User{
		Name:      "Arnold",
		Password:  "321",
		Email:     "arnold@email.com",
		FullName:  "Arnold Schwarzenegger",
		CreatedAt: 1660412960,
	}
)

func TestCache_UserCreate(t *testing.T) {
	testCache := cache{
		mu:     sync.RWMutex{},
		data:   make(map[string]models.User),
		poolCh: make(chan struct{}, 1),
		logger: loggerPkg.NewFatal(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cases := []struct {
		name    string
		user    models.User
		expErr  error
		expUser models.User
		poolCh  func(chan struct{})
	}{
		{
			name:    "success",
			user:    user1,
			expErr:  nil,
			expUser: user1,
			poolCh:  func(_ chan struct{}) {},
		},
		{
			name:    "failed, deadline exceeded",
			user:    user1,
			expErr:  errorsPkg.ErrTimeout,
			expUser: models.User{},
			poolCh: func(ch chan struct{}) {
				ch <- struct{}{}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.poolCh(testCache.poolCh)
			err := testCache.UserCreate(ctx, c.user)
			actualUser := testCache.data[c.user.Name]
			delete(testCache.data, c.user.Name)

			assert.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expUser, actualUser)
		})
	}
}

func TestCache_UserUpdate(t *testing.T) {
	testCache := cache{
		mu:     sync.RWMutex{},
		data:   make(map[string]models.User),
		poolCh: make(chan struct{}, 1),
		logger: loggerPkg.NewFatal(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cases := []struct {
		name    string
		user    models.User
		newUser models.User
		expErr  error
		expUser models.User
		poolCh  func(chan struct{})
	}{
		{
			name:    "success",
			user:    user1,
			newUser: user2,
			expErr:  nil,
			expUser: user2,
			poolCh:  func(_ chan struct{}) {},
		},
		{
			name:    "failed, deadline exceeded",
			user:    user1,
			newUser: user2,
			expErr:  errorsPkg.ErrTimeout,
			expUser: user1,
			poolCh: func(ch chan struct{}) {
				ch <- struct{}{}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testCache.data[c.user.Name] = c.user
			c.poolCh(testCache.poolCh)
			err := testCache.UserUpdate(ctx, c.newUser)
			actualUser := testCache.data[c.newUser.Name]
			delete(testCache.data, c.newUser.Name)

			assert.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expUser, actualUser)
		})
	}
}

func TestCache_UserDelete(t *testing.T) {
	testCache := cache{
		mu:     sync.RWMutex{},
		data:   make(map[string]models.User),
		poolCh: make(chan struct{}, 1),
		logger: loggerPkg.NewFatal(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cases := []struct {
		name    string
		user    models.User
		expErr  error
		expUser models.User
		poolCh  func(chan struct{})
	}{
		{
			name:    "success",
			user:    user1,
			expErr:  nil,
			expUser: models.User{},
			poolCh:  func(_ chan struct{}) {},
		},
		{
			name:    "failed, deadline exceeded",
			user:    user1,
			expErr:  errorsPkg.ErrTimeout,
			expUser: user1,
			poolCh: func(ch chan struct{}) {
				ch <- struct{}{}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testCache.data[c.user.Name] = c.user
			c.poolCh(testCache.poolCh)
			err := testCache.UserDelete(ctx, c.user.Name)
			actualUser := testCache.data[c.user.Name]
			delete(testCache.data, c.user.Name)

			assert.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expUser, actualUser)
		})
	}
}

func TestCache_UserGet(t *testing.T) {
	testCache := cache{
		mu:     sync.RWMutex{},
		data:   make(map[string]models.User),
		poolCh: make(chan struct{}, 1),
		logger: loggerPkg.NewFatal(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cases := []struct {
		name    string
		user    models.User
		expErr  error
		expUser models.User
		poolCh  func(chan struct{})
	}{
		{
			name:    "success",
			user:    user1,
			expErr:  nil,
			expUser: user1,
			poolCh:  func(_ chan struct{}) {},
		},
		{
			name:    "failed, deadline exceeded",
			user:    user1,
			expErr:  errorsPkg.ErrTimeout,
			expUser: models.User{},
			poolCh: func(ch chan struct{}) {
				ch <- struct{}{}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testCache.data[c.user.Name] = c.user
			c.poolCh(testCache.poolCh)
			actualUser, err := testCache.UserGet(ctx, c.user.Name)
			delete(testCache.data, c.user.Name)

			assert.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expUser, actualUser)
		})
	}
}

func TestCache_UserList(t *testing.T) {
	testCache := cache{
		mu:     sync.RWMutex{},
		data:   make(map[string]models.User),
		poolCh: make(chan struct{}, 1),
		logger: loggerPkg.NewFatal(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	testCache.data[user1.Name] = user1
	testCache.data[user3.Name] = user3
	testCache.data[user4.Name] = user4

	cases := []struct {
		name    string
		list    []models.User
		expErr  error
		expList []models.User
		poolCh  func(chan struct{})
		order   bool
		limit   uint64
		offset  uint64
	}{
		{
			name:    "success, asc, all users",
			list:    []models.User{user1, user3, user4},
			expErr:  nil,
			expList: []models.User{user4, user3, user1},
			poolCh:  func(_ chan struct{}) {},
			order:   false,
			limit:   3,
			offset:  0,
		},
		{
			name:    "success, desc, all users",
			list:    []models.User{user1, user3, user4},
			expErr:  nil,
			expList: []models.User{user1, user3, user4},
			poolCh:  func(_ chan struct{}) {},
			order:   true,
			limit:   3,
			offset:  0,
		},
		{
			name:    "success, very big offset",
			list:    []models.User{user1, user3, user4},
			expErr:  nil,
			expList: make([]models.User, 0),
			poolCh:  func(_ chan struct{}) {},
			order:   true,
			limit:   2,
			offset:  2,
		},
		{
			name:    "success, user list less then limit",
			list:    []models.User{user1, user3, user4},
			expErr:  nil,
			expList: []models.User{user4},
			poolCh:  func(_ chan struct{}) {},
			order:   true,
			limit:   2,
			offset:  1,
		},
		{
			name:    "failed, deadline exceeded",
			list:    []models.User{user1, user3, user4},
			expErr:  errorsPkg.ErrTimeout,
			expList: nil,
			poolCh: func(ch chan struct{}) {
				ch <- struct{}{}
			},
			order:  false,
			limit:  2,
			offset: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.poolCh(testCache.poolCh)
			actuaList, err := testCache.UserList(ctx, c.order, c.limit, c.offset)

			assert.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expList, actuaList)
		})
	}
}

func TestCache_Close(t *testing.T) {
	testCache := cache{
		mu:     sync.RWMutex{},
		data:   make(map[string]models.User),
		poolCh: make(chan struct{}, 1),
		logger: loggerPkg.NewFatal(),
	}

	t.Run("success memory clear", func(t *testing.T) {
		testCache.Close()
		_, ok := <-testCache.poolCh

		assert.Nil(t, testCache.data)
		assert.False(t, ok)
	})
}
