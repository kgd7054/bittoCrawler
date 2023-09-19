package blockscraper

import (
	conf "bittoCralwer/config"
	"bittoCralwer/internal/common"
	"bittoCralwer/internal/model"
	"bittoCralwer/internal/protocol/dto"
	"bittoCralwer/internal/service/blockchain/ether/api"
	"context"
	"fmt"
	"time"
)

// TODO common 이나 해당 스크랩함수 수정및 이동
func StartScrapingBlocks(config *conf.Config, model *model.Repositories) {
	server := &api.BlockServer{
		Config:     config,
		Repository: model,
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			blockResponse, err := server.ImportLatestBlock(context.Background())
			if err != nil {
				// TODO fmt 들 log 로 수정
				fmt.Println("error : ", err)
			} else {

				// TODO: block data parsing
				blockNumber, err := common.HexToDecimal(blockResponse.Number)
				if err != nil {
					fmt.Println("err : ", err)
				}

				baseFeePerGas, err := common.HexToDecimal(blockResponse.BaseFeePerGas)
				if err != nil {
					fmt.Println("err : ", err)
				}

				difficulty, err := common.HexToDecimal(blockResponse.Difficulty)
				if err != nil {
					fmt.Println("err : ", err)
				}

				gasLimit, err := common.HexToDecimal(blockResponse.GasLimit)
				if err != nil {
					fmt.Println("err : ", err)
				}

				gasUsed, err := common.HexToDecimal(blockResponse.GasUsed)
				if err != nil {
					fmt.Println("err : ", err)
				}

				size, err := common.HexToDecimal(blockResponse.Size)
				if err != nil {
					fmt.Println("err : ", err)
				}

				totalDifficulty, err := common.HexToDecimal(blockResponse.TotalDifficulty)
				if err != nil {
					fmt.Println("err : ", err)
				}

				blockData := &dto.EthereumBlock{
					BaseFeePerGas:    baseFeePerGas,
					Difficulty:       difficulty,
					ExtraData:        blockResponse.ExtraData,
					GasLimit:         gasLimit,
					GasUsed:          gasUsed,
					Hash:             blockResponse.Hash,
					LogsBloom:        blockResponse.LogsBloom,
					Miner:            blockResponse.Miner,
					MixHash:          blockResponse.MixHash,
					Nonce:            blockResponse.Nonce,
					Number:           blockNumber,
					ParentHash:       blockResponse.ParentHash,
					ReceiptsRoot:     blockResponse.ReceiptsRoot,
					Sha3Uncles:       blockResponse.Sha3Uncles,
					Size:             size,
					StateRoot:        blockResponse.StateRoot,
					Timestamp:        blockResponse.Timestamp,
					TotalDifficulty:  totalDifficulty,
					Transactions:     blockResponse.Transactions,
					TransactionsRoot: blockResponse.TransactionsRoot,
					Withdrawals:      blockResponse.Withdrawals,
					WithdrawalsRoot:  blockResponse.WithdrawalsRoot,
				}

				scopeDB := model.GetScopeDB()
				if scopeDB != nil {
					err := scopeDB.SaveEthereumBlock(blockData)
					if err != nil {
						fmt.Println("err : ", err)
					}
				}

			}
		}
	}
}
