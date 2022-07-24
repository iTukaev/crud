package yaml

import (
	"log"

	"github.com/spf13/viper"

	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
)

type core struct{}

func MustNew() configPkg.Interface {
	return &core{}
}

func (*core) Init() {
	log.Println("init config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic("Config init: ", err)
	}
}

func (*core) BotKey() string {
	return viper.GetString("key")
}
