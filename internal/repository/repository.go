package repository

import (
	"context"
	"fmt"
	"os"

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
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %v", err)
	}

	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) SaveMetricsGas(lastBlock, safeGasPrice, proposeGasPrice, fastGasPrice, suggestBaseFee string) error {
	// Проверяем, существует ли уже запись с таким lastBlock
	var exists bool
	err := r.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM gas_metrics WHERE lastBlock = $1)", lastBlock).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования блока: %v", err)
	}

	if exists {
		fmt.Println("Блок уже существует. Пропускаем вставку.")
		return nil // Возвращаемся, если запись уже существует
	}

	// Вставляем новую запись в таблицу gas_metrics
	query := `INSERT INTO gas_metrics (lastBlock, safeGasPrice, proposeGasPrice, fastGasPrice, suggestBaseFee) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err = r.pool.Exec(context.Background(), query, lastBlock, safeGasPrice, proposeGasPrice, fastGasPrice, suggestBaseFee)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении метрик: %v", err)
	}

	return nil
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}
