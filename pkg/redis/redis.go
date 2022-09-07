package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func New(ctx context.Context, cfg Config) (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       0,
	})

	if err := rds.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rds, nil
}
