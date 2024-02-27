package app

import (
	"swilly-delivery-service/config"
	"swilly-delivery-service/internal/pkg/log"
	redisclient "swilly-delivery-service/internal/pkg/redis"

	"github.com/gomodule/redigo/redis"
)

type Dependency struct {
	Redis *redis.Pool
}

var AppDependency *Dependency

func Bootstrap() error {
	if _, err := config.LoadAndGetConfig(); err != nil {
		return err
	}
	log.SetLogLevel(config.AppConfig.LogLevel)

	AppDependency = &Dependency{
		Redis: redisclient.NewRedisPool(),
	}

	return nil
}
