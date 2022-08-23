package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
	"gitlab.ozon.dev/iTukaev/homework/pkg/logger/emptylog"
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

func TestRepo_UserCreate(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	cases := []struct {
		name   string
		expErr error
	}{
		{
			name:   "success",
			expErr: nil,
		},
		{
			name:   "failed, exec crashed",
			expErr: errorsPkg.ErrUnexpected,
		},
	}
	query := "INSERT INTO users (name,password,email,full_name,created_at) VALUES ($1,$2,$3,$4,$5)"
	args := []interface{}{user.Name, user.Password, user.Email, user.FullName, user.CreatedAt}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.ExpectExec(query).
				WithArgs(args...).
				WillReturnResult(pgxmock.NewResult("INSERT", 1)).
				WillReturnError(c.expErr)

			r := &repo{
				pool:   mock,
				logger: emptylog.New(),
			}
			err = r.UserCreate(context.Background(), user)
			assert.ErrorIs(t, err, c.expErr)
		})
	}
}

func TestRepo_UserUpdate(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	cases := []struct {
		name   string
		expErr error
	}{
		{
			name:   "success",
			expErr: nil,
		},
		{
			name:   "failed, exec crashed",
			expErr: errorsPkg.ErrUnexpected,
		},
	}
	query := "UPDATE users SET password = $1, email = $2, full_name = $3 WHERE name = $4"
	args := []interface{}{user.Password, user.Email, user.FullName, user.Name}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.ExpectExec(query).
				WithArgs(args...).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1)).
				WillReturnError(c.expErr)

			r := &repo{
				pool:   mock,
				logger: emptylog.New(),
			}
			err = r.UserUpdate(context.Background(), user)
			assert.ErrorIs(t, err, c.expErr)
		})
	}
}

func TestRepo_UserDelete(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	cases := []struct {
		name   string
		expErr error
	}{
		{
			name:   "success",
			expErr: nil,
		},
		{
			name:   "failed, exec crashed",
			expErr: errorsPkg.ErrUnexpected,
		},
	}
	query := "DELETE FROM users WHERE name = $1"
	args := []interface{}{user.Name}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mock.ExpectExec(query).
				WithArgs(args...).
				WillReturnResult(pgxmock.NewResult("DELETE", 1)).
				WillReturnError(c.expErr)

			r := &repo{
				pool:   mock,
				logger: emptylog.New(),
			}
			err = r.UserDelete(context.Background(), user.Name)
			assert.ErrorIs(t, err, c.expErr)
		})
	}
}

func TestRepo_UserGet(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	cases := []struct {
		name   string
		err    error
		expErr error
	}{
		{
			name:   "success",
			err:    nil,
			expErr: nil,
		},
		{
			name:   "failed, query crashed",
			err:    errorsPkg.ErrUnexpected,
			expErr: errorsPkg.ErrUnexpected,
		},
		{
			name:   "failed, no data",
			err:    pgx.ErrNoRows,
			expErr: errorsPkg.ErrUserNotFound,
		},
	}
	query := "SELECT name, password, email, full_name, created_at FROM users WHERE name = $1"
	args := []interface{}{user.Name}

	for _, c := range cases {
		rows := pgxmock.NewRows([]string{nameField, passwordField, emailField, fullNameField, createdAtField}).
			AddRow(user.Name, user.Password, user.Email, user.FullName, user.CreatedAt)
		t.Run(c.name, func(t *testing.T) {
			mock.ExpectQuery(query).
				WithArgs(args...).
				WillReturnRows(rows).
				WillReturnError(c.err).
				RowsWillBeClosed()

			r := &repo{
				pool:   mock,
				logger: emptylog.New(),
			}
			_, err = r.UserGet(context.Background(), user.Name)
			assert.ErrorIs(t, err, c.expErr)
		})
	}
}

func TestRepo_UserList(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	cases := []struct {
		name   string
		err    error
		expErr error
	}{
		{
			name:   "success",
			err:    nil,
			expErr: nil,
		},
		{
			name:   "failed, query crashed",
			err:    errorsPkg.ErrUnexpected,
			expErr: errorsPkg.ErrUnexpected,
		},
	}
	order := true
	limit := uint64(2)
	offset := uint64(0)
	query := fmt.Sprintf("SELECT name, password, email, full_name, created_at "+
		"FROM users ORDER BY name DESC LIMIT %d OFFSET %d", limit, offset)

	for _, c := range cases {
		rows := pgxmock.NewRows([]string{nameField, passwordField, emailField, fullNameField, createdAtField}).
			AddRow(user.Name, user.Password, user.Email, user.FullName, user.CreatedAt).
			AddRow(user.Name, user.Password, user.Email, user.FullName, user.CreatedAt)
		t.Run(c.name, func(t *testing.T) {
			mock.ExpectQuery(query).
				WillReturnRows(rows).
				WillReturnError(c.err)

			r := &repo{
				pool:   mock,
				logger: emptylog.New(),
			}
			_, err = r.UserList(context.Background(), order, limit, offset)
			assert.ErrorIs(t, err, c.expErr)
		})
	}
}
