package config

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func Init() error {
	log.Println("init config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "init config")
	}
	return nil
}

func GetApiKey() string {
	return viper.GetString("key")
}
