package metrics

import (
	"math/big"
	"strings"
	"time"

	"eth_metrics2.0/internal/client"
	"eth_metrics2.0/internal/logger"
	"eth_metrics2.0/internal/repository"
)

// CollectorEtherscan собирает метрики из сети Ethereum через Etherscan API и сохраняет их в базу данных.
type CollectorEtherscan struct {
	ethClient *client.EthereumClientEtherscan
	repo      repository.Repository
}

// NewCollectorEtherscan создает новый экземпляр Collector.
func NewCollectorEtherscan(ethClient *client.EthereumClientEtherscan, repo repository.Repository) *CollectorEtherscan {
	return &CollectorEtherscan{
		ethClient: ethClient,
		repo:      repo,
	}
}

type Collector interface {
	CollectAndSave() error
}

// CollectAndSave собирает метрики из сети Ethereum через Etherscan API и сохраняет их в базу данных.
func (c *CollectorEtherscan) CollectAndSave() error {
	// Логируем начало процесса сбора данных
	logger.Logger.WithField("action", "CollectAndSave").Info("Начинаем сбор данных о последнем блоке")

	// Получаем данные о последнем блоке
	blockData, err := c.ethClient.GetBlockData("latest")
	if err != nil {
		logger.Logger.WithError(err).WithField("action", "CollectAndSave").Error("Ошибка при запросе к API")
		return err
	}

	// Проверяем, если в данных есть результат
	if blockData["result"] != nil {
		block := blockData["result"].(map[string]interface{})

		// Логируем номер блока
		blockNumberInt := hexInt(block, "number")

		// Логируем время создания блока
		timestamp := hexInt(block, "timestamp")
		timeCreated := time.Unix(timestamp.Int64(), 0)

		// Логируем количество использованного газа
		gasUsed := hexInt(block, "gasUsed")

		// Получаем транзакции блока
		transactions := block["transactions"].([]interface{})
		lenBlock := len(transactions)

		// Логируем gasPrice для каждой транзакции
		gasprice := metricTransactions(transactions, "gasPrice")

		// Рассчитываем общую комиссию
		totalFees := new(big.Int).Mul(sumSlice(gasprice), gasUsed)

		// Формируем метрики блока
		metricBlock := map[string]interface{}{
			"Transaction_count": lenBlock,
			"Fees":              totalFees,
			"Created_at":        timeCreated,
		}

		// Логируем перед сохранением в базу
		logger.Logger.WithField("block_number", blockNumberInt.Int64()).Info("Сохраняем метрики для блока")
		err = c.repo.SaveMetrics(blockNumberInt.Int64(), "LastBlock", metricBlock, "block_metrics")
		if err != nil {
			logger.Logger.WithError(err).WithField("block_number", blockNumberInt.Int64()).Error("Ошибка при сохранении данных")
			return err
		}

		// Логируем успешное сохранение
		logger.Logger.WithField("block_number", blockNumberInt.Int64()).Info("Метрики успешно сохранены для блока")
	}

	// Логируем завершение процесса
	logger.Logger.WithField("action", "CollectAndSave").Info("Сбор данных завершен.")
	return nil
}

// Достает метрики из транзакции и выводит срез из них
func metricTransactions(transactions []interface{}, metric string) []*big.Int {
	var metrics []*big.Int
	for _, tx := range transactions {
		txData := tx.(map[string]interface{})
		metricBigInt := hexInt(txData, metric)
		metrics = append(metrics, metricBigInt)
	}
	return metrics
}

// hexBigInt конвертирует hex-строку в *big.Int
func hexInt(data map[string]interface{}, metric string) *big.Int {
	metricHex, ok := data[metric].(string)
	if !ok {
		logger.Logger.WithField("metric", metric).Error("Ошибка: поле отсутствует или имеет неверный тип")
	}

	metricHex = strings.TrimPrefix(metricHex, "0x") // Убираем "0x"
	metricBigInt := new(big.Int)
	_, success := metricBigInt.SetString(metricHex, 16) // Парсим 16-ричное число
	if !success {
		logger.Logger.WithField("metric", metric).Error("Ошибка при преобразовании hex строки в big.Int")
	}
	return metricBigInt
}

// Считает сумму слайса
func sumSlice(slice []*big.Int) *big.Int {
	sum := new(big.Int)
	for _, num := range slice {
		sum.Add(sum, num)
	}
	return sum
}
