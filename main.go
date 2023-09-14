package main

import (
	"bittoCralwer/common"
	conf "bittoCralwer/config"
	"bittoCralwer/ether/api"
	ether "bittoCralwer/ether/proto"
	"bittoCralwer/model"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
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

	//model 모듈 선언
	repositories, err := model.NewRepositories(config)
	if err != nil {
		panic(err)
	}
	// TODO DB connection, query
	_ = repositories

	go startScrapingBlocks(config) // Starting the block scraping in a separate goroutine

	// gRPC server setup (unchanged from your code)
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		// TODO log 수정
		fmt.Println("err : ", err)
	}

	grpcServer := grpc.NewServer()
	ether.RegisterBlockServiceServer(grpcServer, &api.BlockServer{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

	// TODO health check
}

// TODO common 이나 해당 스크랩함수 수정및 이동
func startScrapingBlocks(config *conf.Config) {
	server := &api.BlockServer{
		Config: config,
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			blockResponse, err := server.ImportLatestBlock(context.Background(), &ether.ImportBlockRequest{})
			if err != nil {
				// TODO fmt 들 log 로 수정
				fmt.Println("error : ", err)
			} else {
				fmt.Println("block resp : ", blockResponse)

				bn, err := common.HexToDecimal(blockResponse.Result.Number)
				if err != nil {
					fmt.Println("err : ", err)
				}
				_ = bn
				//fmt.Println("block number : ", blockResponse.Result.Number, blockResponse.Result.Hash, blockResponse.Result.Timestamp)
			}
		}
	}
}
