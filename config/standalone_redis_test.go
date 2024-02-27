package config

import (
	"os"
	"testing"
	"time"
)

func TestNewStandaloneRedisConfig(t *testing.T) {
	// setup
	os.Setenv("STANDALONE_REDIS_HOST", "localhost:6379")
	os.Setenv("STANDALONE_REDIS_POOL_SIZE", "10")
	os.Setenv("STANDALONE_REDIS_POOL_TIMEOUT_MS", "100")
	os.Setenv("STANDALONE_REDIS_READ_TIMEOUT_MS", "200")
	os.Setenv("STANDALONE_REDIS_WRITE_TIMEOUT_MS", "300")
	os.Setenv("STANDALONE_REDIS_KEY_PREFIX", "delivery:")

	defer func() {
		// cleanup
		os.Unsetenv("STANDALONE_REDIS_HOST")
		os.Unsetenv("STANDALONE_REDIS_POOL_SIZE")
		os.Unsetenv("STANDALONE_REDIS_POOL_TIMEOUT_MS")
		os.Unsetenv("STANDALONE_REDIS_READ_TIMEOUT_MS")
		os.Unsetenv("STANDALONE_REDIS_WRITE_TIMEOUT_MS")
		os.Unsetenv("STANDALONE_REDIS_KEY_PREFIX")
	}()

	config := newStandaloneRedisConfig()

	// verify
	expectedConfig := &standaloneRedisConfig{
		RedisHost:        "localhost:6379",
		RedisPoolSize:    10,
		RedisPoolTimeout: 100 * time.Millisecond,
	}

	if *config != *expectedConfig {
		t.Errorf("Configuration mismatch. Got: %v, Expected: %v", config, expectedConfig)
	}
}
