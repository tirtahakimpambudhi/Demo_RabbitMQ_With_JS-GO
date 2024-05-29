package config

import (
	internalConfig "go_test_rabbitmq/internal/config"
)

func SetRabbitMQConfig(cfg *internalConfig.Config) {
	cfg.Add("rabbitmq",map[string]any{
		"protocol" : cfg.Env("MESSAGE_BROKER_PROTOCOL","amqp"),
		"host" : cfg.Env("MESSAGE_BROKER_HOST","localhost"),
		"port" : cfg.Env("MESSAGE_BROKER_PORT","15672"),
		"user" : cfg.Env("MESSAGE_BROKER_USER",""),	
		"password" : cfg.Env("MESSAGE_BROKER_PASSWORD",""),
		"vhost" : cfg.Env("MESSAGE_BROKER_VIRTUAL_HOST",""),
		"cloud" : cfg.Env("MESSAGE_BROKER_CLOUD",""),
	})
}