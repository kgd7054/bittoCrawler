package main

import (
	"bittoCralwer/common"
	conf "bittoCralwer/config"
	"bittoCralwer/ether/api"
	"bittoCralwer/model"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configFlag = flag.String("config", "./config/config.toml", "toml file to use for configuration")
	httpFlag   = flag.Int("http", 0, "router http port")
)

func main() {
	config := conf.NewConfig(*configFlag)

	// http
	if *httpFlag != 0 {
		config.Port.Http = *httpFlag
	}

	// model 모듈 선언
	repositories, err := model.NewRepositories(config)
	if err != nil {
		panic(err)
	}
	// TODO DB connection, query

	go startScrapingBlocks(config, repositories) // Starting the block scraping in a separate goroutine

	// TODO health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port.Http), // Use just the port for the HTTP server
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nGracefully shutting down...")
		srv.Shutdown(nil)
		// TODO: Additional cleanup if needed
		os.Exit(0)
	}()

	fmt.Printf("Starting HTTP server on port %d for health checks...\n", config.Port.Http)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Failed to serve HTTP server: %v", err)
	}
}

// TODO common 이나 해당 스크랩함수 수정및 이동
func startScrapingBlocks(config *conf.Config, model *model.Repositories) {
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
				bn, err := common.HexToDecimal(blockResponse.Number)
				if err != nil {
					fmt.Println("err : ", err)
				}
				_ = bn

				fmt.Println("block response : ", blockResponse)

			}
		}
	}
}
