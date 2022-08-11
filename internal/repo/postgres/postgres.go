package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
)

const (
	usersTable = "users"

	nameField      = "name"
	passwordField  = "password"
	emailField     = "email"
	fullNameField  = "full_name"
	createdAtField = "created_at"

	desc = " DESC"
)

func MustNew(ctx context.Context, host, port, user, password, dbname string) repoPkg.Interface {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal("can't connect to database: ", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatal("ping database error: ", err)
	}

	log.Println("With PostgreSQL started")
	return &repo{
		pool: pool,
	}
}

type repo struct {
	pool *pgxpool.Pool
}

func (r *repo) UserCreate(ctx context.Context, user models.User) error {
	query, args, err := squirrel.Insert(usersTable).
		Columns(nameField, passwordField, emailField, fullNameField, createdAtField).
		Values(user.Name, user.Password, user.Email, user.FullName, user.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "postgres UserCreate: to sql")
	}

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return errors.Wrap(err, "postgres UserCreate: insert")
	}

	return nil
}

func (r *repo) UserUpdate(ctx context.Context, user models.User) error {
	builder := squirrel.Update(usersTable).
		Where(squirrel.Eq{
			nameField: user.Name,
		})

	if user.Email != "" {
		builder = builder.Set(emailField, user.Email)
	}
	if user.Password != "" {
		builder = builder.Set(passwordField, user.Password)
	}
	if user.FullName != "" {
		builder = builder.Set(fullNameField, user.FullName)
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(err, "postgres UserUpdate: to sql")
	}

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return errors.Wrap(err, "postgres UserUpdate: update")
	}

	return nil
}

func (r *repo) UserDelete(ctx context.Context, name string) error {
	query, args, err := squirrel.Delete(usersTable).
		Where(squirrel.Eq{
			nameField: name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "postgres UserDelete: to sql")
	}

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return errors.Wrap(err, "postgres UserDelete: delete")
	}

	return nil
}

func (r *repo) UserGet(ctx context.Context, name string) (models.User, error) {
	query, args, err := squirrel.Select(nameField, passwordField, emailField, fullNameField, createdAtField).
		From(usersTable).
		Where(squirrel.Eq{
			nameField: name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return models.User{}, errors.Wrap(err, "postgres UserGet: to sql")
	}
	row := r.pool.QueryRow(ctx, query, args...)
	var user models.User
	if err = row.Scan(&user.Name, &user.Password, &user.Email, &user.FullName, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errorsPkg.ErrUserNotFound
		}
		return models.User{}, errors.Wrap(err, "postgres UserGet: get")
	}

	return user, nil
}

func (r *repo) UserList(ctx context.Context, order bool, limit, offset uint64) ([]models.User, error) {
	var sort string
	if order {
		sort = desc
	}
	query, args, err := squirrel.Select(nameField, passwordField, emailField, fullNameField, createdAtField).
		From(usersTable).
		Limit(limit).
		Offset(offset * limit).
		OrderBy(nameField + sort).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "postgres UserList: to sql")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "postgres UserList: query")
	}

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.Name, &user.Password, &user.Email, &user.FullName, &user.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "postgres UserList: row scan")
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *repo) Close() {
	r.pool.Close()
}
