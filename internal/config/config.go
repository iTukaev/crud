package config

type Interface interface {
	Init()
	BotKey() string
	ServerAddr() string
}
