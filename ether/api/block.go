package api

import (
	conf "bittoCralwer/config"
	ether "bittoCralwer/ether/proto"
	"context"
	"encoding/json"
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

	// Extract details from the blockData
	blockHash := blockData["result"].(map[string]interface{})["hash"].(string)
	blockNumber := blockData["result"].(map[string]interface{})["number"].(string)
	// ... extract other fields as needed

	response := &ether.ImportBlockResponse{
		BlockHash:   blockHash,
		BlockNumber: blockNumber,
		Status:      "imported",
	}

	return response, nil
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
