package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg, err := LoadAndGetConfig()
	if err != nil {
		return
	}
	assert.IsType(t, &Config{}, cfg)
}
