package user

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
	repoMockPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/mock"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

var (
	user = models.User{
		Name:      "Ivan",
		Password:  "123",
		Email:     "ivan@email.com",
		FullName:  "Ivan the Dummy",
		CreatedAt: 1660412940,
	}
)

func Test_Create(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cases := []struct {
		name      string
		user      models.User
		getErr    error
		createErr error
		expErr    error
	}{
		{
			name:      "success",
			user:      user,
			getErr:    errorsPkg.ErrUserNotFound,
			createErr: nil,
			expErr:    nil,
		},
		{
			name:      "failed UserGet unexpected error",
			user:      user,
			getErr:    errorsPkg.ErrUnexpected,
			createErr: nil,
			expErr:    errorsPkg.ErrUnexpected,
		},
		{
			name:      "failed UserGet already exists error",
			user:      user,
			getErr:    nil,
			createErr: nil,
			expErr:    errorsPkg.ErrUserAlreadyExists,
		},
		{
			name:      "failed UserCreate unexpected error",
			user:      user,
			getErr:    errorsPkg.ErrUserNotFound,
			createErr: errorsPkg.ErrUnexpected,
			expErr:    errorsPkg.ErrUnexpected,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := repoMockPkg.NewMockInterface(ctl)
			gomock.InOrder(
				mockRepo.EXPECT().UserGet(gomock.Any(), c.user.Name).
					Return(models.User{}, c.getErr).Times(1),
				mockRepo.EXPECT().UserCreate(gomock.Any(), c.user).
					Return(c.createErr).MaxTimes(1),
			)

			userCtl := New(mockRepo, loggerPkg.NewFatal())
			err := userCtl.Create(context.Background(), c.user)
			assert.ErrorIs(t, err, c.expErr)
		})
	}
}

func Test_Update(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cases := []struct {
		name      string
		user      models.User
		getErr    error
		updateErr error
		expErr    error
	}{
		{
			name:      "success",
			user:      user,
			getErr:    nil,
			updateErr: nil,
			expErr:    nil,
		},
		{
			name:      "failed UserGet unexpected error",
			user:      user,
			getErr:    errorsPkg.ErrUnexpected,
			updateErr: nil,
			expErr:    errorsPkg.ErrUnexpected,
		},
		{
			name:      "failed UserUpdate unexpected error",
			user:      user,
			getErr:    nil,
			updateErr: errorsPkg.ErrUnexpected,
			expErr:    errorsPkg.ErrUnexpected,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := repoMockPkg.NewMockInterface(ctl)
			gomock.InOrder(
				mockRepo.EXPECT().UserGet(gomock.Any(), c.user.Name).
					Return(models.User{}, c.getErr).Times(1),
				mockRepo.EXPECT().UserUpdate(gomock.Any(), c.user).
					Return(c.updateErr).MaxTimes(1),
			)

			userCtl := New(mockRepo, loggerPkg.NewFatal())
			err := userCtl.Update(context.Background(), c.user)
			assert.ErrorIs(t, err, c.expErr)
		})
	}
}

func Test_Delete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cases := []struct {
		name      string
		user      string
		getErr    error
		deleteErr error
		expErr    error
	}{
		{
			name:      "success",
			user:      user.Name,
			getErr:    nil,
			deleteErr: nil,
			expErr:    nil,
		},
		{
			name:      "failed UserGet unexpected error",
			user:      user.Name,
			getErr:    errorsPkg.ErrUnexpected,
			deleteErr: nil,
			expErr:    errorsPkg.ErrUnexpected,
		},
		{
			name:      "failed UserDelete unexpected error",
			user:      user.Name,
			getErr:    nil,
			deleteErr: errorsPkg.ErrUnexpected,
			expErr:    errorsPkg.ErrUnexpected,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := repoMockPkg.NewMockInterface(ctl)
			gomock.InOrder(
				mockRepo.EXPECT().UserGet(gomock.Any(), c.user).
					Return(models.User{}, c.getErr).Times(1),
				mockRepo.EXPECT().UserDelete(gomock.Any(), c.user).
					Return(c.deleteErr).MaxTimes(1),
			)

			userCtl := New(mockRepo, loggerPkg.NewFatal())
			err := userCtl.Delete(context.Background(), c.user)
			assert.ErrorIs(t, err, c.expErr)
		})
	}
}

func Test_Get(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cases := []struct {
		name    string
		user    string
		getErr  error
		expErr  error
		expUser models.User
	}{
		{
			name:    "success",
			user:    user.Name,
			getErr:  nil,
			expErr:  nil,
			expUser: user,
		},
		{
			name:    "failed UserGet unexpected error",
			user:    user.Name,
			getErr:  errorsPkg.ErrUnexpected,
			expErr:  errorsPkg.ErrUnexpected,
			expUser: models.User{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := repoMockPkg.NewMockInterface(ctl)
			gomock.InOrder(
				mockRepo.EXPECT().UserGet(gomock.Any(), c.user).
					Return(c.expUser, c.getErr).Times(1),
			)

			userCtl := New(mockRepo, loggerPkg.NewFatal())
			expUser, err := userCtl.Get(context.Background(), c.user)
			assert.ErrorIs(t, err, c.expErr)
			assert.Equal(t, expUser, c.expUser)
		})
	}
}

func Test_List(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cases := []struct {
		name    string
		listErr error
		expErr  error
		expList []models.User
	}{
		{
			name:    "success",
			listErr: nil,
			expErr:  nil,
			expList: []models.User{user},
		},
		{
			name:    "failed UserList unexpected error",
			listErr: errorsPkg.ErrUnexpected,
			expErr:  errorsPkg.ErrUnexpected,
			expList: make([]models.User, 0),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := repoMockPkg.NewMockInterface(ctl)
			gomock.InOrder(
				mockRepo.EXPECT().UserList(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(c.expList, c.listErr).Times(1),
			)

			userCtl := New(mockRepo, loggerPkg.NewFatal())
			expList, err := userCtl.List(context.Background(), true, 1, 1)
			assert.ErrorIs(t, err, c.expErr)
			assert.Equal(t, expList, c.expList)
		})
	}
}
