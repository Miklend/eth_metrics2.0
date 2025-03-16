package config

import (
	"os"

	"eth_metrics2.0/internal/logger"
	"github.com/joho/godotenv"
)

// Config хранит настройки приложения
type Config struct {
	PostgresUser      string
	PostgresPassword  string
	PostgresHost      string
	PostgresPort      string
	PostgresDB        string
	ETHERSCAN_API_KEY string
}

// Load загружает конфигурацию из .env или переменных окружения
func Load() *Config {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		logger.Logger.Warn("Файл .env не найден, используются переменные окружения")
	} else {
		logger.Logger.Info("Файл .env успешно загружен")
	}

	cfg := &Config{
		PostgresUser:      getEnv("POSTGRES_USER"),
		PostgresPassword:  getEnv("POSTGRES_PASSWORD"),
		PostgresHost:      getEnv("POSTGRES_HOST"),
		PostgresPort:      getEnv("POSTGRES_PORT"),
		PostgresDB:        getEnv("POSTGRES_DB"),
		ETHERSCAN_API_KEY: getEnv("ETHERSCAN_API_KEY"),
	}

	logger.Logger.Info("Конфигурация успешно загружена")

	return cfg
}

// getEnv получает переменную окружения и логирует, если её нет
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Logger.WithField("key", key).Warn("Переменная окружения не задана")
	}
	return value
}
