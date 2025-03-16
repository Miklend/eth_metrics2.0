package main

import (
	"fmt"
	"time"

	"eth_metrics2.0/internal/client"
	"eth_metrics2.0/internal/config"
	"eth_metrics2.0/internal/logger"
	"eth_metrics2.0/internal/metrics"
	"eth_metrics2.0/internal/repository"
	"github.com/cenkalti/backoff/v4"
)

func main() {
	// Инициализация логгера
	logger.Init()

	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к PostgreSQL с повторными попытками
	var repo repository.Repository
	retryBackoff := backoff.NewExponentialBackOff()
	retryBackoff.MaxElapsedTime = time.Minute * 5 // Ограничение на общее время попыток

	err := backoff.Retry(func() error {
		var err error
		repo, err = repository.NewPostgresRepository()
		if err != nil {
			return fmt.Errorf("не удалось подключиться к PostgreSQL: %v", err)
		}
		return nil
	}, retryBackoff)

	if err != nil {
		logger.Logger.WithError(err).Fatal("Не удалось подключиться к PostgreSQL после нескольких попыток")
	}
	defer repo.Close()

	// Создание клиента для работы с Etherscan API
	ethClient := client.NewEthereumClientEtherscan(cfg.ETHERSCAN_API_KEY)
	logger.Logger.Info("Создан клиент для работы с Etherscan API")

	// Создание сборщика метрик
	collectorGas := metrics.NewCollectorEtherscan(ethClient, repo)

	// Канал для передачи результатов и ошибок
	resultCh := make(chan string)
	errorCh := make(chan error)

	// Интервал для сбора метрик
	interval := 8 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Бесконечный цикл для сбора метрик каждые 8 секунд
	go func() {
		for range ticker.C {
			// Сбор и сохранение метрик с повторными попытками
			err := backoff.Retry(func() error {
				return collectorGas.CollectAndSaveGas()
			}, retryBackoff)

			if err != nil {
				errorCh <- fmt.Errorf("Ошибка при сборе метрик: %v", err)
			} else {
				resultCh <- "Метрики успешно собраны и сохранены."
			}
		}
	}()

	// Обработка результатов и ошибок
	for {
		select {
		case result := <-resultCh:
			logger.Logger.Info(result)
		case err := <-errorCh:
			logger.Logger.WithError(err).Error("Ошибка при сборе метрик")
		}
	}
}
