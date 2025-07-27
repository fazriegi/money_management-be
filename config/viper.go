package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var viperInstance *viper.Viper

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("../../")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viperInstance = config

	return config
}

func GetConfigString(key string) string {
	return viperInstance.GetString(key)
}

func GetConfigInt(key string) int {
	return viperInstance.GetInt(key)
}
