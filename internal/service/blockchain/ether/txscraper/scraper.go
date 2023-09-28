package txscraper

import (
	conf "bittoCralwer/config"
	"bittoCralwer/internal/common"
	"bittoCralwer/internal/model"
	"bittoCralwer/internal/service/blockchain/ether/api"
	"log"
	"time"
)

const (
	interval = 5 * time.Second
)

func StartScrapingTxs(config *conf.Config, model *model.Repositories) {
	log.Println("starting tx scraping")
	server := &api.TransactionServer{
		Config:     config,
		Repository: model,
	}
	scopeDB := model.GetScopeDB()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		transactions, err := scopeDB.GetEthereumTx()
		if err != nil {
			common.Warn("StartScrapingTxs", "GetEthereumTx", err)
			continue
		}
		for _, tx := range transactions {

			txInfo, err := server.GetTransactionDetailData(tx.Transaction)
			if err != nil {
				common.Error("StartScrapingTxs", "GetTransactionDetailData", err, "hash", tx)
				continue
			}

			err = scopeDB.InsertEthereumTxDetailInfo(*txInfo)
			if err != nil {
				common.Error("StartScrapingTxs", "InsertEthereumTxDetailInfo", err)
				continue
			}

			err = scopeDB.MarkAsProcessedTxData(tx.Id)
			if err != nil {
				common.Error("StartScrapingTxs", "MarkAsProcessedTxData", err)
				continue
			}
		}
	}
}
