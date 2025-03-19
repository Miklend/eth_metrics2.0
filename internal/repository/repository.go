package repository

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	SaveMetrics(primaryKeyValue interface{}, primaryKeyField string, metrics map[string]interface{}, nameTable string) error
	Close()
}

// Метод для сохранения метрик газа
func (r *PostgresRepository) SaveMetrics(primaryKeyValue interface{}, primaryKeyField string, metrics map[string]interface{}, nameTable string) error {
	logger.Logger.WithField("primaryKeyValue", primaryKeyValue).Info("Проверка существования метрик")

	// Проверяем, существует ли уже запись с таким primaryKeyValue
	var exists bool
	queryCheck := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)", nameTable, primaryKeyField)
	err := r.pool.QueryRow(context.Background(), queryCheck, primaryKeyValue).Scan(&exists)
	if err != nil {
		logger.Logger.WithError(err).WithField("primaryKeyValue", primaryKeyValue).Error("Ошибка при проверке существования записи")
		return fmt.Errorf("ошибка при проверке существования записи: %v", err)
	}

	// Динамическое формирование SQL-запросов
	fields := []string{}
	placeholders := []string{}
	values := []interface{}{primaryKeyValue} // Используем primaryKeyValue вместо lastBlock
	updateParts := []string{}

	i := 2 // Начинаем с $2, так как $1 уже занят primaryKeyValue
	for key, value := range metrics {
		fields = append(fields, key)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		updateParts = append(updateParts, fmt.Sprintf("%s = $%d", key, i))
		values = append(values, value)
		i++
	}

	if exists {
		// Обновляем данные
		updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $1", nameTable, strings.Join(updateParts, ", "), primaryKeyField)
		_, err = r.pool.Exec(context.Background(), updateQuery, values...)
		if err != nil {
			logger.Logger.WithError(err).WithField("primaryKeyValue", primaryKeyValue).Error("Ошибка при обновлении метрик")
			return fmt.Errorf("ошибка при обновлении метрик для записи: %v", err)
		}
		logger.Logger.WithField("primaryKeyValue", primaryKeyValue).Info("Метрики успешно обновлены")
		return nil
	}

	// Вставляем новую запись
	insertQuery := fmt.Sprintf(
		"INSERT INTO %s (%s, %s) VALUES ($1, %s)",
		nameTable,
		primaryKeyField, // Теперь используем primaryKeyField
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err = r.pool.Exec(context.Background(), insertQuery, values...)
	if err != nil {
		logger.Logger.WithError(err).WithField("primaryKeyValue", primaryKeyValue).Error("Ошибка при сохранении метрик")
		return fmt.Errorf("ошибка при сохранении метрик: %v", err)
	}

	return nil
}

// Закрытие соединения
func (r *PostgresRepository) Close() {
	logger.Logger.Info("Закрытие соединения с базой данных")
	r.pool.Close()
}
