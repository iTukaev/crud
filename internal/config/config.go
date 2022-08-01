package config

import pgModels "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres/models"

type Interface interface {
	Init()
	BotKey() string
	GRPCAddr() string
	HTTPAddr() string
	PGConfig() pgModels.Config
}
