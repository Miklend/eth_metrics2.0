package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EthereumClientEtherscan хранит методы для работы с Etherscan API.
type EthereumClientEtherscan struct {
	ApiKey string
}

// NewEthereumClientEtherscan создаёт новый клиент для работы с Etherscan API.
func NewEthereumClientEtherscan(apiKey string) *EthereumClientEtherscan {
	return &EthereumClientEtherscan{ApiKey: apiKey}
}

func (e *EthereumClientEtherscan) GetBlockData(blockNumber string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s", blockNumber, e.ApiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
