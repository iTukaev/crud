package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"

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

type PgxPool interface {
	pgxtype.Querier
	Close()
}

func New(pool *pgxpool.Pool, logger *zap.SugaredLogger) repoPkg.Interface {
	logger.Infoln("With PostgreSQL started")
	return &repo{
		pool:   pool,
		logger: logger,
	}
}

func NewPostgres(ctx context.Context, host, port, user, password, dbname string, logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	logger.Debugln("PostgreSQL connection", psqlConn)

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to database: %v\n", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database error: %v\n", err)
	}
	return pool, nil
}

type repo struct {
	pool   PgxPool
	logger *zap.SugaredLogger
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
	r.logger.Debugln("UserCreate", query, args)

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return errors.Wrap(err, "postgres UserCreate: insert")
	}

	return nil
}

func (r *repo) UserUpdate(ctx context.Context, user models.User) error {
	query, args, err := squirrel.Update(usersTable).
		Set(passwordField, user.Password).
		Set(emailField, user.Email).
		Set(fullNameField, user.FullName).
		Where(squirrel.Eq{
			nameField: user.Name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "postgres UserUpdate: to sql")
	}
	r.logger.Debugln("UserUpdate", query, args)

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
	r.logger.Debugln("UserDelete", query, args)

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
	r.logger.Debugln("UserGet", query, args)

	row := r.pool.QueryRow(ctx, query, args...)
	var user models.User
	if err = row.Scan(&user.Name, &user.Password, &user.Email, &user.FullName, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errorsPkg.ErrUserNotFound
		}
		return models.User{}, errors.Wrap(err, "postgres UserGet: get")
	}
	r.logger.Debugln("UserGet", user.String())

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
	r.logger.Debugln("UserList", query, args)

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
	r.logger.Debugln("UserList", users)

	return users, nil
}

func (r *repo) Close() {
	r.pool.Close()
	r.logger.Infoln("PostgreSQL connection closed")
}
