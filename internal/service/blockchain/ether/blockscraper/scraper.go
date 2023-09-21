package blockscraper

import (
	conf "bittoCralwer/config"
	"bittoCralwer/internal/common"
	"bittoCralwer/internal/model"
	"bittoCralwer/internal/protocol/dto"
	"bittoCralwer/internal/service/blockchain/ether/api"
	"context"
	"time"
)

func StartScrapingBlocks(config *conf.Config, model *model.Repositories) {
	server := &api.BlockServer{
		Config:     config,
		Repository: model,
	}
	redisDB := model.GetRedisDB()

	err := SyncWithAlchemy(model, server)
	if err != nil {
		common.Error("StartScrapingBlocks", "SyncWithAlchemy", err)
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			blockResponse, err := server.ImportLatestBlock(context.Background())
			if err != nil {
				// TODO common -> log로 수정
				common.Error("StartScrapingBlocks", "ImportLatestBlock", err)
				continue

			} else {

				// TODO: block data parsing
				blockNumber, err := common.HexToDecimal(blockResponse.Number)
				if err != nil {
					common.Error("StartScrapingBlocks", "HexToDecimal", err, "blockNumber", blockNumber)
					continue
				}

				baseFeePerGas, err := common.HexToDecimal(blockResponse.BaseFeePerGas)
				if err != nil {
					common.Error("StartScrapingBlocks", "HexToDecimal", err, "baseFeePerGas", baseFeePerGas)
					continue
				}

				difficulty, err := common.HexToDecimal(blockResponse.Difficulty)
				if err != nil {
					common.Error("StartScrapingBlocks", "HexToDecimal", err, "difficulty", difficulty)
					continue
				}

				gasLimit, err := common.HexToDecimal(blockResponse.GasLimit)
				if err != nil {
					common.Error("StartScrapingBlocks", "HexToDecimal", err, "gasLimit", gasLimit)
					continue
				}

				gasUsed, err := common.HexToDecimal(blockResponse.GasUsed)
				if err != nil {
					common.Error("StartScrapingBlocks", "HexToDecimal", err, "gasUsed", gasUsed)
					continue
				}

				size, err := common.HexToDecimal(blockResponse.Size)
				if err != nil {
					common.Error("StartScrapingBlocks", "HexToDecimal", err, "size", size)
					continue
				}

				totalDifficulty, err := common.HexToDecimal(blockResponse.TotalDifficulty)
				if err != nil {
					common.Error("StartScrapingBlocks", "HexToDecimal", err, "totalDifficulty", totalDifficulty)
					continue
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

				err = redisDB.SetCache("last_block_number", blockNumber)
				if err != nil {
					common.Error("StartScrapingBlocks", "SetCache", err)
					continue
				}

				scopeDB := model.GetScopeDB()
				if scopeDB != nil {
					err := scopeDB.SaveEthereumBlock(blockData)
					if err != nil {
						common.Warn("StartScrapingBlocks", "SaveEthereumBlock", err)
						continue
					}
				}

				common.Info("StartScrapingBlocks", "blockNumber", blockData.Number)

			}
		}
	}
}

func SyncWithAlchemy(model *model.Repositories, server *api.BlockServer) error {

	redisDB := model.GetRedisDB()
	lastProcessedBlock, err := redisDB.GetCache("last_block_number")
	if err != nil {
		return err
	}

	blockData, err := server.ImportLatestBlock(context.TODO())
	if err != nil {
		return err
	}

	blockDataNumber, err := common.HexToDecimal(blockData.Number)
	if err != nil {
		return err
	}

	processedBlockNumber, err := common.ConvertStringToInt(lastProcessedBlock)
	if err != nil {
		return err
	}

	latestBlockNumber, err := common.ConvertStringToInt(blockDataNumber)
	if err != nil {
		return err
	}

	for i := processedBlockNumber + 1; i <= latestBlockNumber; i++ {
		missedBlock, err := server.GetBlockByNumber(i)
		if err != nil {
			return err
		}

		// TODO: 재사용된 함수 및 DB save 함수화 및 모듈화
		blockNumber, err := common.HexToDecimal(missedBlock.Number)
		if err != nil {
			common.Error("StartScrapingBlocks", "HexToDecimal", err, "blockNumber", blockNumber)
			continue
		}

		baseFeePerGas, err := common.HexToDecimal(missedBlock.BaseFeePerGas)
		if err != nil {
			common.Error("StartScrapingBlocks", "HexToDecimal", err, "baseFeePerGas", baseFeePerGas)
			continue
		}

		difficulty, err := common.HexToDecimal(missedBlock.Difficulty)
		if err != nil {
			common.Error("StartScrapingBlocks", "HexToDecimal", err, "difficulty", difficulty)
			continue
		}

		gasLimit, err := common.HexToDecimal(missedBlock.GasLimit)
		if err != nil {
			common.Error("StartScrapingBlocks", "HexToDecimal", err, "gasLimit", gasLimit)
			continue
		}

		gasUsed, err := common.HexToDecimal(missedBlock.GasUsed)
		if err != nil {
			common.Error("StartScrapingBlocks", "HexToDecimal", err, "gasUsed", gasUsed)
			continue
		}

		size, err := common.HexToDecimal(missedBlock.Size)
		if err != nil {
			common.Error("StartScrapingBlocks", "HexToDecimal", err, "size", size)
			continue
		}

		totalDifficulty, err := common.HexToDecimal(missedBlock.TotalDifficulty)
		if err != nil {
			common.Error("StartScrapingBlocks", "HexToDecimal", err, "totalDifficulty", totalDifficulty)
			continue
		}

		blockData := &dto.EthereumBlock{
			BaseFeePerGas:    baseFeePerGas,
			Difficulty:       difficulty,
			ExtraData:        missedBlock.ExtraData,
			GasLimit:         gasLimit,
			GasUsed:          gasUsed,
			Hash:             missedBlock.Hash,
			LogsBloom:        missedBlock.LogsBloom,
			Miner:            missedBlock.Miner,
			MixHash:          missedBlock.MixHash,
			Nonce:            missedBlock.Nonce,
			Number:           blockNumber,
			ParentHash:       missedBlock.ParentHash,
			ReceiptsRoot:     missedBlock.ReceiptsRoot,
			Sha3Uncles:       missedBlock.Sha3Uncles,
			Size:             size,
			StateRoot:        missedBlock.StateRoot,
			Timestamp:        missedBlock.Timestamp,
			TotalDifficulty:  totalDifficulty,
			Transactions:     missedBlock.Transactions,
			TransactionsRoot: missedBlock.TransactionsRoot,
			Withdrawals:      missedBlock.Withdrawals,
			WithdrawalsRoot:  missedBlock.WithdrawalsRoot,
		}

		scopeDB := model.GetScopeDB()
		if scopeDB != nil {
			err := scopeDB.SaveEthereumBlock(blockData)
			if err != nil {
				common.Warn("StartScrapingBlocks", "SaveEthereumBlock", err)
				continue
			}
		}
	}

	return nil
}
