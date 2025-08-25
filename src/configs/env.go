package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DB_NAME                   string
	DB_HOST                   string
	DB_USER                   string
	DB_PASSWORD               string
	DB_PORT                   int
	DB_SSLMODE                string
	JWT_EXPIRATION_IN_SECONDS int64
	JWT_SECRET                string
}

var EnvConfigs = initConfig()

func initConfig() EnvConfig {
	godotenv.Load()

	return EnvConfig{
		DB_NAME:                   getEnv("DB_NAME", ""),
		DB_HOST:                   getEnv("DB_HOST", "localhost"),
		DB_USER:                   getEnv("DB_USER", ""),
		DB_PASSWORD:               getEnv("DB_PASSWORD", ""),
		DB_PORT:                   int(getEnvAsInt("DB_PORT", 5432)),
		DB_SSLMODE:                getEnv("DB_SSLMODE", ""),
		JWT_EXPIRATION_IN_SECONDS: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600),
		JWT_SECRET:                getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
