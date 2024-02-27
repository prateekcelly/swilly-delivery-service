package config

import (
	"time"
)

type standaloneRedisConfig struct {
	RedisHost        string
	RedisPoolSize    int
	RedisPoolTimeout time.Duration
}

func newStandaloneRedisConfig() *standaloneRedisConfig {
	return &standaloneRedisConfig{
		RedisHost:        getStringOrPanic("STANDALONE_REDIS_HOST"),
		RedisPoolSize:    getIntOrPanic("STANDALONE_REDIS_POOL_SIZE"),
		RedisPoolTimeout: time.Millisecond * time.Duration(getIntOrPanic("STANDALONE_REDIS_POOL_TIMEOUT_MS")),
	}
}
