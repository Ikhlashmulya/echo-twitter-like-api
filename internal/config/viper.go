package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	viper := viper.New()

	viper.AddConfigPath("..")
	viper.AddConfigPath("..")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("error load configuration : %v", err.Error()))
	}

	return viper
}