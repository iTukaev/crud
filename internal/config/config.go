package config

import pgModels "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres/models"

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
}

type Data interface {
	RepoAddr() string
	PGConfig() pgModels.Config
	Local() bool
	WorkersCount() int
}
