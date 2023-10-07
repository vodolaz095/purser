package config

import "os"

func loadFromEnvironment(v *string, key string) {
	if fromEnv := os.Getenv(key); fromEnv != "" {
		*v = fromEnv
	}
}
