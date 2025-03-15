package main

import (
	"fmt"
	"log"
	"time"

	"eth_metrics2.0/internal/client"
	"eth_metrics2.0/internal/config"
	"eth_metrics2.0/internal/metrics"
	"eth_metrics2.0/internal/repository"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к PostgreSQL
	repo, err := repository.NewPostgresRepository()
	if err != nil {
		log.Fatalf("Не удалось подключиться к PostgreSQL: %v", err)
	}
	defer repo.Close()

	// Создание клиента для работы с Etherscan API
	ethClient := client.NewEthereumClientEtherscan(cfg.ETHERSCAN_API_KEY)

	// Создание сборщика метрик
	collectorGas := metrics.NewCollectorEtherscan(ethClient, repo)

	// Интервал в 5 секунд
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Бесконечный цикл для сбора метрик каждые 5 секунд
	for range ticker.C {
		// Сбор и сохранение метрик
		if err := collectorGas.CoolectAndSaveGas(); err != nil {
			log.Fatalf("Ошибка при сборе метрик: %v", err)
		} else {
			fmt.Println("Метрики успешно собраны и сохранены.", time.Now())
		}
	}
}
