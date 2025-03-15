package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost      string
	PostgresPort      string
	PostgresUser      string
	PostgresPassword  string
	PostgresDB        string
	ETHERSCAN_API_KEY string
	LogLevel          string
}

func Load() *Config {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}

	return &Config{
		PostgresHost:      getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:      getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:      getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword:  getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:        getEnv("POSTGRES_DB", "eth_metrics2"),
		ETHERSCAN_API_KEY: getEnv("ETHERSCAN_API_KEY", ""),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
