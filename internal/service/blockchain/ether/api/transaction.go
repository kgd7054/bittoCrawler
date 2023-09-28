package api

import (
	conf "bittoCralwer/config"
	"bittoCralwer/internal/common"
	"bittoCralwer/internal/model"
	"bittoCralwer/internal/protocol/dto"
	"context"
	"encoding/hex"
	"fmt"
	ethercommon "github.com/ethereum/go-ethereum/common"
	ethertypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type TransactionServer struct {
	Config     *conf.Config
	Repository *model.Repositories
}

func (t *TransactionServer) GetTransactionDetailData(hash string) (*dto.TransactionDetail, error) {
	client, err := ethclient.Dial(alchemyAPIURL + t.Config.Alchemy.ApiKey)
	if err != nil {
		return nil, err
	}
	strHash := ethercommon.HexToHash(hash)

	tx, isPending, err := client.TransactionByHash(context.Background(), strHash)
	if err != nil {
		return nil, err
	}
	if isPending {
		common.Error("GetTransactionDetailData", "TransactionByHash", "still pending", "tx", tx)
		return nil, fmt.Errorf("still pendgin tx")
	}

	receipt, err := client.TransactionReceipt(context.Background(), strHash)
	if err != nil {
		return nil, err
	}

	blockNumber := receipt.BlockNumber
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, err
	}

	inputData := tx.Data()
	hexInputData := hex.EncodeToString(inputData)
	method := ""
	if len(tx.Data()) == 0 {
		method = ""
	} else {
		method = hex.EncodeToString(inputData[:4])
	}

	from := ethercommon.Address{}
	to := ethercommon.Address{}
	gasPrice := big.NewInt(0)
	gasUsed := receipt.GasUsed
	baseFee := block.BaseFee()
	txFee := big.NewInt(0)
	maxFee := big.NewInt(0)
	maxPriorityFee := big.NewInt(0)
	totalBurntFee := big.NewInt(0)
	txSavings := big.NewInt(0)
	switch tx.Type() {
	case ethertypes.LegacyTxType:
		signer := ethertypes.NewEIP155Signer(tx.ChainId())
		from, err = ethertypes.Sender(signer, tx)
		gasPrice = tx.GasPrice()
		txFee = new(big.Int).Mul(big.NewInt(int64(gasUsed)), gasPrice)

	case ethertypes.DynamicFeeTxType:
		from, err = ethertypes.Sender(ethertypes.NewLondonSigner(tx.ChainId()), tx)
		if baseFee == nil {
			baseFee = tx.GasPrice()
		}
		gasPrice = new(big.Int).Add(baseFee, tx.EffectiveGasTipValue(baseFee))
		txFee = new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasUsed))
		maxFee = tx.GasFeeCap()
		if baseFee != nil {
			//maxPriorityFee = new(big.Int).Sub(tx.GasTipCap(), baseFee)
			//maxPriorityFee = math.BigMax(maxPriorityFee, big.NewInt(0))
			maxPriorityFee = tx.GasTipCap()
		}
		totalBurntFee = new(big.Int).Mul(baseFee, new(big.Int).SetUint64(gasUsed))
		actualTip := common.MinBigInt(maxPriorityFee, new(big.Int).Sub(maxFee, baseFee))
		actualTotalFee := new(big.Int).Add(baseFee, actualTip)
		txSavings = new(big.Int).Sub(maxFee, actualTotalFee)
		txSavings = new(big.Int).Mul(txSavings, new(big.Int).SetUint64(gasUsed))

	default:
		common.Error("GetTransactionDetailData", "tx.Type", "invalid tx type")
	}

	if tx.To() == nil {
		to = receipt.ContractAddress
	}

	result := dto.TransactionDetail{
		Hash:               hash,
		Status:             receipt.Status,
		BlockNumber:        blockNumber,
		Timestamp:          block.Time(),
		Method:             "0x" + method,
		From:               from.String(),
		To:                 to.String(),
		Value:              tx.Value(),
		TransactionFee:     txFee,
		GasPrice:           gasPrice,
		GasLimit:           tx.Gas(),
		GasUsed:            gasUsed, // usage by txn
		GasFeesBase:        baseFee,
		GasFeesMax:         maxFee,
		GasFeesMaxPriority: maxPriorityFee,
		Burnt:              totalBurntFee,
		TxnSavingsFees:     txSavings,
		Nonce:              tx.Nonce(),
		TxType:             tx.Type(),
		InputData:          "0x" + hexInputData,
	}

	return &result, nil
}
