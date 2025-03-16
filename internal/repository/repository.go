package repository

import (
	"context"
	"fmt"
	"os"

	"eth_metrics2.0/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository() (*PostgresRepository, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Logger.WithError(err).Error("Не удалось подключиться к PostgreSQL")
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %v", err)
	}

	logger.Logger.Info("Успешное подключение к PostgreSQL")
	return &PostgresRepository{pool: pool}, nil
}

// Интерфейс репозитория
type Repository interface {
	SaveMetricsGas(lastBlock, safeGasPrice, proposeGasPrice, fastGasPrice, suggestBaseFee string) error
	Close()
}

// Метод для сохранения метрик газа
func (r *PostgresRepository) SaveMetricsGas(lastBlock, safeGasPrice, proposeGasPrice, fastGasPrice, suggestBaseFee string) error {
	logger.Logger.WithField("lastBlock", lastBlock).Info("Проверка существования метрик газа")

	// Проверяем, существует ли уже запись с таким lastBlock
	var exists bool
	err := r.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM gas_metrics WHERE lastBlock = $1)", lastBlock).Scan(&exists)
	if err != nil {
		logger.Logger.WithError(err).WithField("lastBlock", lastBlock).Error("Ошибка при проверке существования блока")
		return fmt.Errorf("ошибка при проверке существования блока: %v", err)
	}

	if exists {
		// Обновляем данные, если блок уже есть
		updateQuery := `UPDATE gas_metrics SET safeGasPrice = $2, proposeGasPrice = $3, fastGasPrice = $4, suggestBaseFee = $5 
						WHERE lastBlock = $1`
		_, err := r.pool.Exec(context.Background(), updateQuery, lastBlock, safeGasPrice, proposeGasPrice, fastGasPrice, suggestBaseFee)
		if err != nil {
			logger.Logger.WithError(err).WithField("lastBlock", lastBlock).Error("Ошибка при обновлении метрик")
			return fmt.Errorf("ошибка при обновлении метрик для блока: %v", err)
		}

		logger.Logger.WithField("lastBlock", lastBlock).Info("Метрики успешно обновлены")
		return nil
	}

	// Вставляем новую запись
	insertQuery := `INSERT INTO gas_metrics (lastBlock, safeGasPrice, proposeGasPrice, fastGasPrice, suggestBaseFee) 
                    VALUES ($1, $2, $3, $4, $5)`
	_, err = r.pool.Exec(context.Background(), insertQuery, lastBlock, safeGasPrice, proposeGasPrice, fastGasPrice, suggestBaseFee)
	if err != nil {
		logger.Logger.WithError(err).WithField("lastBlock", lastBlock).Error("Ошибка при сохранении метрик")
		return fmt.Errorf("ошибка при сохранении метрик: %v", err)
	}

	return nil
}

// Закрытие соединения
func (r *PostgresRepository) Close() {
	logger.Logger.Info("Закрытие соединения с базой данных")
	r.pool.Close()
}
