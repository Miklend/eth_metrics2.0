package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"eth_metrics2.0/internal/client"
	"eth_metrics2.0/internal/repository"
)

// Collector собирает метрики из сети Ethereum через Etherscan API и сохраняет их в базу данных.
type CollectorEtherscan struct {
	ethClient *client.EthereumClientEtherscan
	repo      *repository.PostgresRepository
}

// NewCollector создает новый экземпляр Collector.
func NewCollectorEtherscan(ethClient *client.EthereumClientEtherscan, repo *repository.PostgresRepository) *CollectorEtherscan {
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

// CoolectAndSaveGas собирает метрики по газу(медленная, средная, быстрая, базовая комиссия) из сети Ethereum и сохраняет их в базу данных.
func (c *CollectorEtherscan) CoolectAndSaveGas() error {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=gastracker&action=gasoracle&apikey=%s", c.ethClient.ApiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка при запросе к API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	var gasStats GasStatsResponse
	err = json.Unmarshal(body, &gasStats)
	if err != nil {
		log.Fatalf("Ошибка при разборе JSON: %v", err)
	}

	err = c.repo.SaveMetricsGas(gasStats.Result.LastBlock, gasStats.Result.SafeGasPrice, gasStats.Result.ProposeGasPrice, gasStats.Result.FastGasPrice, gasStats.Result.SuggestBaseFee)
	if err != nil {
		log.Fatalf("Ошибка при сохранении метрик газа: %v", err)
	}
	return nil
}
