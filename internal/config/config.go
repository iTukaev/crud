package config

type Interface interface {
	Init()
	BotKey() string
	GRPCAddr() string
	HTTPAddr() string
}
