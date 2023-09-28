package dao

import (
	"bittoCralwer/internal/common"
	"bittoCralwer/internal/protocol/dto"
	"database/sql"
)

type EthereumTxDAO struct {
	DB *sql.DB
}

func NewEthereumTxDAO(db *sql.DB) *EthereumTxDAO {
	return &EthereumTxDAO{DB: db}
}

func (dao *EthereumTxDAO) GetEtherTx() ([]dto.TransactionsInfo, error) {

	rows, err := dao.DB.Query("SELECT * FROM transactions WHERE processed = FALSE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]dto.TransactionsInfo, 0)
	for rows.Next() {
		tx := dto.TransactionsInfo{}
		err := rows.Scan(&tx.Id, &tx.BlockId, &tx.Transaction, &tx.Processed)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	// TODO: mark transactions as processed
	if len(transactions) == 0 {
		common.Info("GetEtherTx", "transactions", "length is 0")
	}

	return transactions, nil

}

func (dao *EthereumTxDAO) InsertEthTx(detail *dto.TransactionDetail) error {
	query := `
    INSERT INTO tx_detail_info (
        hash, status, block_number, timestamp, method,
        from_address, to_address, value, transaction_fee, gas_price,
        gas_used, gas_limit, nonce, tx_type, gas_fees_base,
        gas_fees_max, gas_fees_max_priority, burnt, txn_savings_fees, input_data
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
  `

	_, err := dao.DB.Exec(
		query,
		detail.Hash,
		detail.Status,
		detail.BlockNumber.String(), // assuming BlockNumber is of type *big.Int
		detail.Timestamp,
		detail.Method,
		detail.From,
		detail.To,
		detail.Value.String(),
		detail.TransactionFee.String(),
		detail.GasPrice.String(),
		detail.GasUsed,
		detail.GasLimit,
		detail.Nonce,
		detail.TxType,
		detail.GasFeesBase.String(),
		detail.GasFeesMax.String(),
		detail.GasFeesMaxPriority.String(),
		detail.Burnt.String(),
		detail.TxnSavingsFees.String(),
		detail.InputData,
	)

	return err
}

func (dao *EthereumTxDAO) MarkAsProcessed(id int) error {
	_, err := dao.DB.Exec("UPDATE transactions SET processed = TRUE WHERE id = ?", id)
	return err
}
