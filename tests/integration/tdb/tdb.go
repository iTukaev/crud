//go:build integration
// +build integration

package tdb

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	postgresPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

const (
	Host     = "localhost"
	Port     = "5432"
	User     = "user"
	Password = "password"
	DBName   = "candy_shop"
)

func NewTestDB(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := postgresPkg.NewPostgres(ctx, Host, Port, User, Password, DBName, loggerPkg.NewFatal())
	if err != nil {
		return nil, err
	}
	return pool, nil
}
