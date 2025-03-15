package client

// EthereumClientEtherscan хранит методы для работы с Etherscan API.
type EthereumClientEtherscan struct {
	ApiKey string
}

// NewEthereumClientEtherscan создаёт новый клиент для работы с Etherscan API.
func NewEthereumClientEtherscan(apiKey string) *EthereumClientEtherscan {
	return &EthereumClientEtherscan{ApiKey: apiKey}
}
