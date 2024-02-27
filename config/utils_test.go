package config

import (
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetStringWithDefault(t *testing.T) {
	key := "STRING_DEFAULT"
	os.Unsetenv(key)
	assert.Equal(t, getStringWithDefault(key, "DEFAULT"), "DEFAULT")
}

func TestGetIntWithDefault(t *testing.T) {
	key := "INT_DEFAULT"
	os.Unsetenv(key)
	assert.Equal(t, getIntWithDefault(key, 100), 100)
}

func TestGetBoolWithDefault(t *testing.T) {
	key := "BOOL_DEFAULT"
	os.Unsetenv(key)
	assert.Equal(t, getBoolWithDefault(key, false), false)
}
