package config

import (
	pgModels "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres/models"
	redisPkg "gitlab.ozon.dev/iTukaev/homework/pkg/redis"
)

type Interface interface {
	Transport
	Data
	ExternalServices
}

type ExternalServices interface {
	LogLevel() string
	Brokers() []string
	JService() string
	JHost() string
}

type Transport interface {
	BotKey() string
	GRPCAddr() string
	GRPCDataAddr() string
	HTTPAddr() string
	HTTPDataAddr() string
}

type Data interface {
	PGConfig() pgModels.Config
	Local() bool
	WorkersCount() int
	RedisConfig() redisPkg.Config
}
