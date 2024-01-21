package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger(config *viper.Viper) *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	logLevel, err := logrus.ParseLevel(config.GetString("log.level"))
	if err != nil {
		panic(fmt.Errorf("error parsing log level: %v", err))
	}

	logger.SetLevel(logLevel)

	return logger
}