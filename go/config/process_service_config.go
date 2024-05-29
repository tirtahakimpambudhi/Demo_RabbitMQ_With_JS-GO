package config

import (
	internalConfig"go_test_rabbitmq/internal/config"
)

func SetProcessServiceConfig(cfg *internalConfig.Config) {
	cfg.Add("process_service",map[string]any{
		"routing_key" : cfg.Env("ROUTING_KEY_PROCESS_SERVICE",""),
		"exchange" : cfg.Env("EXCHANGE_PROCESS_SERVICE",""),
		"queue" : cfg.Env("QUEUE_PROCESS_SERVICE",""), 
	})
}