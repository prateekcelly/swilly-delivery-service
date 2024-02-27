package config

import (
	"github.com/spf13/viper"
)

const (
	serviceName = "swilly-delivery-service"
)

type Config struct {
	HTTPServerPort        string
	LogLevel              string
	ServiceName           string
	WorkerEnabled         bool
	DirectoryPath         string
	JobName               string
	StandaloneRedisConfig *standaloneRedisConfig
}

var AppConfig *Config

func LoadAndGetConfig() (*Config, error) {

	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")
	viper.AddConfigPath("../../../../")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viper.AutomaticEnv()
	AppConfig = &Config{
		HTTPServerPort:        getStringWithDefault("HTTP_SERVER_PORT", "8080"),
		ServiceName:           serviceName,
		LogLevel:              getStringWithDefault("LOG_LEVEL", "info"),
		WorkerEnabled:         getBoolWithDefault("WORKER_ENABLED", true),
		DirectoryPath:         getStringOrPanic("DIRECTORY_PATH"),
		JobName:               getStringWithDefault("JOB_NAME", "send_message"),
		StandaloneRedisConfig: newStandaloneRedisConfig(),
	}
	return AppConfig, nil
}
