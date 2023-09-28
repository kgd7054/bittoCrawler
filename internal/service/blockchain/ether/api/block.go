package api

import (
	conf "bittoCralwer/config"
	"bittoCralwer/internal/common"
	"bittoCralwer/internal/model"
	"bittoCralwer/internal/protocol/dto"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ethercommon "github.com/ethereum/go-ethereum/common"
	ethertypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// TODO: 공통 상수로 빼기
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
		Withdrawals:      block.Withdrawals,
		WithdrawalsRoot:  block.WithdrawalsRoot,
	}, nil
}

func (s *BlockServer) GetBlockByNumber(blockNumber int) (*dto.EthereumBlock, error) {

	blockNumberHex, err := common.DecimalStringToHex(strconv.Itoa(blockNumber))
	if err != nil {
		return nil, err
	}

	payload := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "eth_getBlockByNumber",
		"params": ["%s", false]
	}`, blockNumberHex)

	req, err := http.NewRequest("POST", alchemyAPIURL+s.Config.Alchemy.ApiKey, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	alchemyResponse := AlchemyResponse{}
	err = json.Unmarshal(body, &alchemyResponse)
	if err != nil {
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
		Withdrawals:      block.Withdrawals,
		WithdrawalsRoot:  block.WithdrawalsRoot,
	}, nil

}

func (s *BlockServer) getLatestBlockFromAlchemy() ([]byte, error) {

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

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *BlockServer) GetTransactionInfo() {
	//======
	hash := "0x6863329f0665ab3e8057081d97fe5b04d7e8270ab4a7917aba0188c739772917"
	client, err := ethclient.Dial(alchemyAPIURL + s.Config.Alchemy.ApiKey)
	if err != nil {
		fmt.Println("err : ", err)
	}
	strHash := ethercommon.HexToHash(hash)

	fmt.Println("client : ", client)
	fmt.Println("hash : ", strHash)

	tx, _, err := client.TransactionByHash(context.TODO(), strHash)
	if err != nil {
		fmt.Println("err : ", err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), strHash)
	if err != nil {
		fmt.Println("err : ", err)
	}
	fmt.Println("receipt : ", receipt)
	blockNumber := receipt.BlockNumber
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		fmt.Println("block : ", block)
	}

	timestamp := block.Time()
	fmt.Println("time stamp : ", timestamp)

	inputData := tx.Data()
	//fmt.Println("input data : ", inputData)

	hexInputData := hex.EncodeToString(inputData)
	fmt.Printf("hex input data : 0x%s", hexInputData)

	method := hex.EncodeToString(inputData[:4])
	fmt.Printf("method : 0x%s", method)

	from := ethercommon.Address{}
	switch tx.Type() {
	case ethertypes.LegacyTxType:
		from, err = ethertypes.Sender(ethertypes.HomesteadSigner{}, tx)
	case ethertypes.DynamicFeeTxType:
		from, err = ethertypes.Sender(ethertypes.NewLondonSigner(tx.ChainId()), tx)
	default:
		fmt.Println("not")
	}

	if err != nil {
		fmt.Println("err : ", err)
	}

	fmt.Println("from : ", from)
	fmt.Println("to : ", tx.To())

	//mars, err := tx.MarshalJSON()
	//if err != nil {
	//	fmt.Println("err : ", err)
	//}
	//fmt.Println("marshal : ", mars)

	// 매주 화요일 저녁 10시 위클리, 주간 회의같은 것
	//======
}
