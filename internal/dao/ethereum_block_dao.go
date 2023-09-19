package dao

import (
	"bittoCralwer/internal/protocol/dto"
	"database/sql"
	"fmt"
)

type EthereumBlockDAO struct {
	DB *sql.DB
}

func NewEthereumBlockDAO(db *sql.DB) *EthereumBlockDAO {
	return &EthereumBlockDAO{DB: db}
}

// Save saves the EthereumBlock into the database.
func (dao *EthereumBlockDAO) Save(block *dto.EthereumBlock) error {
	// Convert and save the EthereumBlock object into your MySQL database.
	// Use DAO.DB to execute the SQL query.
	// Return any error that might occur.

	tx, err := dao.DB.Begin()
	if err != nil {
		return err
	}

	blockInsert := `
    INSERT INTO blocks (
        BaseFeePerGas, Difficulty, ExtraData, GasLimit, GasUsed, Hash, LogsBloom, 
        Miner, MixHash, Nonce, Number, ParentHash, ReceiptsRoot, Sha3Uncles, 
        Size, StateRoot, Timestamp, TotalDifficulty, TransactionsRoot, WithdrawalsRoot
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

	// Get the ID of the last inserted block
	blockID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("block id : ", blockID)

	// Insert transactions
	for _, transaction := range block.Transactions {
		_, err = tx.Exec("INSERT INTO transactions (block_id, transaction) VALUES (?, ?)", blockID, transaction)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert uncles
	for _, uncle := range block.Uncles {
		_, err = tx.Exec("INSERT INTO uncles (block_id, uncle) VALUES (?, ?)", blockID, uncle)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert withdrawals
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
