package config

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Config struct {
	vip *viper.Viper
}

func NewConfig() (*Config, error) {
	vip := viper.New()
	vip.AddConfigPath(".")
	vip.AddConfigPath("../")
	vip.AddConfigPath("../../")
	vip.SetConfigName(".env")
	vip.SetConfigType("env")
	vip.AutomaticEnv() // Automatically read environment variables

	err := vip.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Log the error but do not return, allowing env variables to be used
			fmt.Println("Config file not found; using environment variables only")
		} else {
			// If there was another error (e.g., malformed file), return it
			return nil, err
		}
	}

	return &Config{vip: vip}, nil
}

func (app *Config) Add(name string, configuration any) {
	app.vip.Set(name, configuration)
}

func (app *Config) Env(envName string, defaultValue ...any) any {
	value := app.Get(envName, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return nil
	}

	return value
}

// Get config from Config.
func (app *Config) Get(path string, defaultValue ...any) any {
	if !app.vip.IsSet(path) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	return app.vip.Get(path)
}

// GetString Get string type config from Config.
func (app *Config) GetString(path string, defaultValue ...any) string {
	value := cast.ToString(app.Get(path, defaultValue...))
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(string)
		}

		return ""
	}

	return value
}

// GetInt Get int type config from Config.
func (app *Config) GetInt(path string, defaultValue ...any) int {
	value := app.Get(path, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(int)
		}

		return 0
	}

	return cast.ToInt(value)
}

// GetBool Get bool type config from Config.
func (app *Config) GetBool(path string, defaultValue ...any) bool {
	value := app.Get(path, defaultValue...)
	if cast.ToString(value) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0].(bool)
		}

		return false
	}

	return cast.ToBool(value)
}
