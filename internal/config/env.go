package config

import (
	"fmt"
	"os"
)

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		var parsed int
		_, err := fmt.Sscanf(val, "%d", &parsed)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
