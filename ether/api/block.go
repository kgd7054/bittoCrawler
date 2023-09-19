package api

import (
	conf "bittoCralwer/config"
	"bittoCralwer/model"
	"bittoCralwer/protocol/dto"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const alchemyAPIURL = "https://eth-sepolia.g.alchemy.com/v2/"

type BlockServer struct {
	Config     *conf.Config
	Repository *model.Repositories
}

type AlchemyResponse struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      int               `json:"id"`
	Result  dto.EthereumBlock `json:"result"`
}

func (s *BlockServer) ImportLatestBlock(ctx context.Context) (*dto.EthereumBlock, error) {

	blockData, err := s.getLatestBlockFromAlchemy()
	if err != nil {
		return nil, err
	}

	var alchemyResponse AlchemyResponse
	if err := json.Unmarshal(blockData, &alchemyResponse); err != nil {
		return nil, err
	}

	block := alchemyResponse.Result

	return &dto.EthereumBlock{
		BaseFeePerGas:    block.BaseFeePerGas,
		Difficulty:       block.Difficulty,
		ExtraData:        block.ExtraData,
		GasLimit:         block.GasLimit,
		GasUsed:          block.GasUsed,
		Hash:             block.Hash,
		LogsBloom:        block.LogsBloom,
		Miner:            block.Miner,
		MixHash:          block.MixHash,
		Nonce:            block.Nonce,
		Number:           block.Number,
		ParentHash:       block.ParentHash,
		ReceiptsRoot:     block.ReceiptsRoot,
		Sha3Uncles:       block.Sha3Uncles,
		Size:             block.Size,
		StateRoot:        block.StateRoot,
		Timestamp:        block.Timestamp,
		TotalDifficulty:  block.TotalDifficulty,
		Transactions:     block.Transactions,
		TransactionsRoot: block.TransactionsRoot,
		Withdrawal:       block.Withdrawal,
		WithdrawalsRoot:  block.WithdrawalsRoot,
	}, nil
}

func (s *BlockServer) getLatestBlockFromAlchemy() ([]byte, error) {
	// Create a new request
	// Define the request payload
	payload := `{
       "jsonrpc": "2.0",
       "id": 1,
       "method": "eth_getBlockByNumber",
       "params": ["latest", false]
   }`

	req, err := http.NewRequest("POST", alchemyAPIURL+s.Config.Alchemy.ApiKey, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response data
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
