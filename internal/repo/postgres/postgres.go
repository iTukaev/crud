package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
)

const (
	usersTable = "users"

	nameField      = "name"
	passwordField  = "password"
	emailField     = "email"
	fullNameField  = "full_name"
	createdAtField = "created_at"
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

	return &repo{
		pool: pool,
	}
}

type repo struct {
	pool *pgxpool.Pool
}

func (r *repo) UserCreate(ctx context.Context, user models.User) error {
	query, args, err := squirrel.Insert(usersTable).
		Columns(nameField, passwordField, emailField, fullNameField).
		Values(user.Name, user.Password, user.Email, user.FullName).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	fmt.Println(query)
	if err != nil {
		return fmt.Errorf("UserCreate: to sql: %w", err)
	}
	row := r.pool.QueryRow(ctx, query, args...)
	if err = row.Scan(&user.Name); err != nil {
		return fmt.Errorf("UserCreate: insert: %w", err)
	}

	return nil
}

func (r *repo) UserUpdate(ctx context.Context, user models.User) error {
	query, args, err := squirrel.Update(usersTable).
		Where(squirrel.Eq{
			nameField: user.Name,
		}).
		Set(passwordField, user.Password).
		Set(emailField, user.Email).
		Set(fullNameField, user.FullName).
		ToSql()
	fmt.Println(query)
	if err != nil {
		return fmt.Errorf("UserUpdate: to sql: %w", err)
	}
	row := r.pool.QueryRow(ctx, query, args...)
	if err = row.Scan(&user.Name); err != nil {
		return fmt.Errorf("UserUpdate: update: %w", err)
	}

	return nil
}

func (r *repo) UserDelete(ctx context.Context, name string) error {
	query, args, err := squirrel.Delete(usersTable).
		Where(squirrel.Eq{
			nameField: name,
		}).ToSql()
	fmt.Println(query)
	if err != nil {
		return fmt.Errorf("UserDelete: to sql: %w", err)
	}
	row := r.pool.QueryRow(ctx, query, args...)
	if err = row.Scan(&name); err != nil {
		return fmt.Errorf("UserDelete: delete: %w", err)
	}

	return nil
}

func (r *repo) UserGet(ctx context.Context, name string) (models.User, error) {
	query, args, err := squirrel.Select(nameField, passwordField, emailField, fullNameField).
		From(usersTable).
		Where(squirrel.Eq{
			nameField: name,
		}).ToSql()
	fmt.Println(query)
	if err != nil {
		return models.User{}, fmt.Errorf("UserGet: to sql: %w", err)
	}
	row := r.pool.QueryRow(ctx, query, args...)
	var user models.User
	if err = row.Scan(&user.Name, &user.Password, &user.Email, &user.FullName); err != nil {
		return models.User{}, fmt.Errorf("UserGet: get: %w", err)
	}

	return user, nil
}

func (r *repo) UserList(ctx context.Context, order bool, limit, offset uint32) ([]models.User, error) {
	//query, args, err := squirrel.Select(nameField, passwordField, emailField, fullNameField).
	//	From(usersTable).
	//	ToSql()
	//fmt.Println(query)
	//if err != nil {
	//	return nil, fmt.Errorf("Repository.UserGet: to sql: %w", err)
	//}
	//row := r.pool.
	//var user models.User
	//if err = row.Scan(&user.Name, &user.Password, &user.Email, &user.FullName); err != nil {
	//	return nil, fmt.Errorf("Repository.UserGet: get: %w", err)
	//}

	return nil, nil
}

func (r *repo) Close() {
	r.pool.Close()
}
