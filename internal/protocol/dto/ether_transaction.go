package dto

import "math/big"

type TransactionsInfo struct {
	Id          int    `json:"id"`
	BlockId     int    `json:"block_id"`
	Transaction string `json:"transaction"`
	Processed   bool   `json:"processed"`
}

type TransactionDetail struct {
	Hash               string
	Status             uint64
	BlockNumber        *big.Int
	Timestamp          uint64
	Method             string
	From               string
	To                 string
	Value              *big.Int
	TransactionFee     *big.Int
	GasPrice           *big.Int
	InputData          string
	GasUsed            uint64
	Nonce              uint64
	TxType             uint8
	GasLimit           uint64
	GasFeesBase        *big.Int
	GasFeesMax         *big.Int
	GasFeesMaxPriority *big.Int
	Burnt              *big.Int
	TxnSavingsFees     *big.Int
	// TODO: gas limit & usage by txn, gas fees, burnt & txn savings fees, other attribuutes
}
