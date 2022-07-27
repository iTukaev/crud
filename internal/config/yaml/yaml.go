package yaml

import (
	"log"

	"github.com/spf13/viper"

	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
)

type config struct{}

func MustNew() configPkg.Interface {
	return &config{}
}

func (config) Init() {
	log.Println("init config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Config init: ", err)
	}
}

func (config) BotKey() string {
	return viper.GetString("key")
}

func (config) ServerAddr() string {
	return viper.GetString("addr")
}
