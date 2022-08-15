package yaml

import (
	"log"

	"github.com/spf13/viper"

	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	pgModels "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres/models"
)

type config struct{}

func MustNew() configPkg.Interface {
	log.Println("Init config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config init: %v\n", err)
	}
	return &config{}
}

func (config) BotKey() string {
	return viper.GetString("key")
}

func (config) GRPCAddr() string {
	return viper.GetString("grpc")
}

func (config) HTTPAddr() string {
	return viper.GetString("http")
}

func (config) RepoAddr() string {
	return viper.GetString("repo")
}

func (config) PGConfig() pgModels.Config {
	var pg pgModels.Config
	if err := viper.UnmarshalKey("pg", &pg); err != nil {
		log.Fatalf("Postgres config unmarshal error: %v\n", err)
	}
	return pg
}

func (config) Local() bool {
	return viper.GetBool("local")
}

func (config) WorkersCount() int {
	return viper.GetInt("workers")
}
