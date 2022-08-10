package config

import pgModels "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres/models"

type Interface interface {
	BotKey() string
	GRPCAddr() string
	HTTPAddr() string
	Local() bool
	WorkersCount() int
	PGConfig() pgModels.Config
}
