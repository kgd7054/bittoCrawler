package dao

import (
	"bittoCralwer/internal/protocol/dto"
	"database/sql"
)

type EthereumBlockDAO struct {
	DB *sql.DB
}

func NewEthereumBlockDAO(db *sql.DB) *EthereumBlockDAO {
	return &EthereumBlockDAO{DB: db}
}

// Save saves the EthereumBlock into the database.
func (dao *EthereumBlockDAO) SaveEtherBlock(block *dto.EthereumBlock) error {

	tx, err := dao.DB.Begin()
	if err != nil {
		return err
	}

	blockInsert := `
    INSERT INTO blocks (
        base_fee_per_gas, difficulty, extra_data, gas_limit, gas_used, hash, logs_bloom, 
        miner, mix_hash, nonce, number, parent_hash, receipts_root, sha_3_uncles, 
        size, state_root, timestamp, total_difficulty, transactions_root, withdrawals_root
    )
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(blockInsert, block.BaseFeePerGas, block.Difficulty, block.ExtraData, block.GasLimit, block.GasUsed,
		block.Hash, block.LogsBloom, block.Miner, block.MixHash, block.Nonce, block.Number, block.ParentHash, block.ReceiptsRoot,
		block.Sha3Uncles, block.Size, block.StateRoot, block.Timestamp, block.TotalDifficulty, block.TransactionsRoot,
		block.WithdrawalsRoot) // Removed the withdrawal fields here since they'll be inserted into a separate table
	if err != nil {
		tx.Rollback()
		return err
	}

	blockID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// TODO: table save 작은 단위로 나누기
	for _, transaction := range block.Transactions {
		_, err = tx.Exec("INSERT INTO transactions (block_id, transaction, processed) VALUES (?, ?, ?)", blockID, transaction, false)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, uncle := range block.Uncles {
		_, err = tx.Exec("INSERT INTO uncles (block_id, uncle) VALUES (?, ?)", blockID, uncle)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, withdrawal := range block.Withdrawals {
		_, err = tx.Exec("INSERT INTO withdrawals (block_id, address, amount, withdrawal_index, validator_index) VALUES (?, ?, ?, ?, ?)", blockID, withdrawal.Address, withdrawal.Amount, withdrawal.Index, withdrawal.ValidatorIndex)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Get retrieves the EthereumBlock from the database.
func (dao *EthereumBlockDAO) Get(blockHash string) (*dto.EthereumBlock, error) {
	// Use DAO.DB to fetch the block details using the given blockHash.
	// Convert the result into a EthereumBlock object.
	// Return the EthereumBlock object and any error that might occur.
	return nil, nil
}
