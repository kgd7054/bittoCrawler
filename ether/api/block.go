package api

import (
	conf "bittoCralwer/config"
	ether "bittoCralwer/ether/proto"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const alchemyAPIURL = "https://eth-sepolia.g.alchemy.com/v2/"

type BlockServer struct {
	ether.UnimplementedBlockServiceServer
	Config *conf.Config
}

func (s *BlockServer) ImportLatestBlock(ctx context.Context, req *ether.ImportBlockRequest) (*ether.ImportBlockResponse, error) {
	// Implement the logic to import the latest Ethereum block using the Alchemy Sepolia API.
	// Construct and return a proto.ImportBlockResponse object populated with the block data.
	blockData, err := s.getLatestBlockFromAlchemy()
	if err != nil {
		return nil, err
	}

	result, ok := blockData["result"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to cast result to map[string]interface{}")
	}

	block := &ether.Block{
		BaseFeePerGas:    result["baseFeePerGas"].(string),
		Difficulty:       result["difficulty"].(string),
		ExtraData:        result["extraData"].(string),
		GasLimit:         result["gasLimit"].(string),
		GasUsed:          result["gasUsed"].(string),
		Hash:             result["hash"].(string),
		LogsBloom:        result["logsBloom"].(string),
		Miner:            result["miner"].(string),
		MixHash:          result["mixHash"].(string),
		Nonce:            result["nonce"].(string),
		Number:           result["number"].(string),
		ParentHash:       result["parentHash"].(string),
		ReceiptsRoot:     result["receiptsRoot"].(string),
		Sha3Uncles:       result["sha3Uncles"].(string),
		Size:             result["size"].(string),
		StateRoot:        result["stateRoot"].(string),
		Timestamp:        result["timestamp"].(string),
		TotalDifficulty:  result["totalDifficulty"].(string),
		TransactionsRoot: result["transactionsRoot"].(string),
		WithdrawalsRoot:  result["withdrawalsRoot"].(string),
	}

	if txs, ok := result["transactions"].([]interface{}); ok {
		for _, tx := range txs {
			block.Transactions = append(block.Transactions, tx.(string))
		}
	}

	if uncles, ok := result["uncles"].([]interface{}); ok {
		for _, uncle := range uncles {
			block.Uncles = append(block.Uncles, uncle.(string))
		}
	}

	//blockNumber := ethercommon.HexToAddress(block.Number)
	//fmt.Println("block number : ", blockNumber)

	//bn, err := common.HexToDecimal(block.Number)
	//if err != nil {
	//	fmt.Println("err : ", err)
	//}
	//fmt.Println("bn : ", bn)

	return &ether.ImportBlockResponse{
		Result: block,
	}, nil
}

func (s *BlockServer) getLatestBlockFromAlchemy() (map[string]interface{}, error) {
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

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
