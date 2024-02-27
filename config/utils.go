package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func checkKey(key string) {
	if !viper.IsSet(key) && os.Getenv(key) == "" {
		log.Fatalf("%s key is not set", key)
	}
}

func getStringWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntWithDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

func getStringOrPanic(key string) string {
	checkKey(key)
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	return value
}

func getIntOrPanic(key string) int {
	checkKey(key)
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("%s key is not Integer value", key)
	}
	return intVal
}

func getBoolWithDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		value = viper.GetString(key)
	}

	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}
