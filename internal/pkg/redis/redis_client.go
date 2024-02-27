package redisclient

import (
	"swilly-delivery-service/config"

	"github.com/gomodule/redigo/redis"
)

func NewRedisPool() *redis.Pool {
	redisConfig := config.AppConfig.StandaloneRedisConfig
	return &redis.Pool{
		MaxActive:   redisConfig.RedisPoolSize,
		MaxIdle:     redisConfig.RedisPoolSize,
		IdleTimeout: redisConfig.RedisPoolTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisConfig.RedisHost)
		},
	}
}
