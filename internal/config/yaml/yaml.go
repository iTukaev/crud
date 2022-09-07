package yaml

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	pgModels "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres/models"
	redisPkg "gitlab.ozon.dev/iTukaev/homework/pkg/redis"
)

type config struct{}

func New() (configPkg.Interface, error) {
	log.Println("Init config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "config init")
	}
	return &config{}, nil
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

func (config) GRPCDataAddr() string {
	return viper.GetString("grpc_data")
}

func (config) HTTPDataAddr() string {
	return viper.GetString("http_data")
}

func (config) LogLevel() string {
	return viper.GetString("log")
}

func (config) PGConfig() pgModels.Config {
	var pg pgModels.Config
	if err := viper.UnmarshalKey("pg", &pg); err != nil {
		log.Fatalf("Postgres config unmarshal error: %v\n", err)
	}
	return pg
}

func (config) RedisConfig() redisPkg.Config {
	var cfg redisPkg.Config
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		log.Fatalf("Postgres config unmarshal error: %v\n", err)
	}
	return cfg
}

func (config) Local() bool {
	return viper.GetBool("local")
}

func (config) WorkersCount() int {
	return viper.GetInt("workers")
}

func (config) Brokers() []string {
	return viper.GetStringSlice("brokers")
}

func (config) JService() string {
	return viper.GetString("jaeger.service")
}

func (config) JHost() string {
	return viper.GetString("jaeger.host")
}
