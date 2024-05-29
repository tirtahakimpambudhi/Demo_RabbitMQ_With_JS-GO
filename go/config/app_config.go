package config

import (
	"fmt"
	cfg "go_test_rabbitmq/internal/config"
	"go_test_rabbitmq/internal/logger"
)

func NewAppConfig() *cfg.Config {
	log := logger.NewLogger()
	appConfig , err := cfg.NewConfig()
	if err != nil {
		fmt.Println(err.Error())
		log.Errorln(err.Error())
	}
	// SET CONFIG
	appConfig.Add("app",map[string]any{
		"env" : appConfig.Env("APP_ENV","testing"),
		"secret_key" : appConfig.Env("SECRET_KEY",""),
	})
	SetRabbitMQConfig(appConfig)
	SetProcessServiceConfig(appConfig)
	return appConfig
}