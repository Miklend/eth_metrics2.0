package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

type GasStatsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		LastBlock       string `json:"LastBlock"`
		SafeGasPrice    string `json:"SafeGasPrice"`
		ProposeGasPrice string `json:"ProposeGasPrice"`
		FastGasPrice    string `json:"FastGasPrice"`
		SuggestBaseFee  string `json:"suggestBaseFee"`
		GasUsedRatio    string `json:"gasUsedRatio"`
	} `json:"result"`
}

// CollectAndSaveGas собирает метрики по газу (медленная, средняя, быстрая, базовая комиссия) из сети Ethereum и сохраняет их в базу данных.
func (c *CollectorEtherscan) CollectAndSaveGas() error {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=%s", c.ethClient.ApiKey)
	logger.Logger.WithField("url", url).Info("Отправка запроса к API Etherscan")

	resp, err := http.Get(url)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка при запросе к API")
		return fmt.Errorf("ошибка при запросе к API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка при чтении ответа от API")
		return fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	var gasStats GasStatsResponse
	err = json.Unmarshal(body, &gasStats)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка при разборе JSON")
		return fmt.Errorf("ошибка при разборе JSON: %v", err)
	}

	logger.Logger.WithFields(map[string]interface{}{
		"lastBlock":       gasStats.Result.LastBlock,
		"safeGasPrice":    gasStats.Result.SafeGasPrice,
		"proposeGasPrice": gasStats.Result.ProposeGasPrice,
		"fastGasPrice":    gasStats.Result.FastGasPrice,
		"suggestBaseFee":  gasStats.Result.SuggestBaseFee,
	}).Info("Успешно получены данные о газе")

	err = c.repo.SaveMetricsGas(
		gasStats.Result.LastBlock,
		gasStats.Result.SafeGasPrice,
		gasStats.Result.ProposeGasPrice,
		gasStats.Result.FastGasPrice,
		gasStats.Result.SuggestBaseFee,
	)
	if err != nil {
		logger.Logger.WithError(err).Error("Ошибка при сохранении метрик газа в базу")
		return fmt.Errorf("ошибка при сохранении метрик газа: %v", err)
	}

	return nil
}
