package main

import (
	"fmt"
	"log"
	"time"

	"eth_metrics2.0/internal/client"
	"eth_metrics2.0/internal/config"
	"eth_metrics2.0/internal/logger"
	"eth_metrics2.0/internal/metrics"
	"eth_metrics2.0/internal/repository"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация логгера
	logger.Init()

	// Подключение к PostgreSQL
	var repo repository.Repository
	repo, err := repository.NewPostgresRepository()
	if err != nil {
		log.Fatalf("Не удалось подключиться к PostgreSQL: %v", err)
	}
	defer repo.Close()

	// Создание клиента для работы с Etherscan API
	ethClient := client.NewEthereumClientEtherscan(cfg.ETHERSCAN_API_KEY)

	// Создание сборщика метрик
	collectorGas := metrics.NewCollectorEtherscan(ethClient, repo)

	// Интервал в 8 секунд, можно заменить на конфигурируемую переменную
	interval := 8 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Бесконечный цикл для сбора метрик каждые 5 секунд
	for range ticker.C {
		// Сбор и сохранение метрик
		if err := collectorGas.CollectAndSaveGas(); err != nil {
			log.Printf("Ошибка при сборе метрик: %v", err)
		} else {
			fmt.Println("Метрики успешно собраны и сохранены.", time.Now())
		}
	}
}
