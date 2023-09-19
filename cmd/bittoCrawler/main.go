package main

import (
	conf "bittoCralwer/config"
	"bittoCralwer/internal/model"
	"bittoCralwer/internal/server"
	"bittoCralwer/internal/service/blockchain/ether/blockscraper"
	"flag"
	"fmt"
	"sync"
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

	repositories, err := model.InitializeRepositories(config)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		blockscraper.StartScrapingBlocks(config, repositories)
	}()

	srv := server.StartServer(config.Port.Http)

	// This function can also accept the waitgroup if you want to track the server goroutine
	server.HandleGracefulShutdown(srv)

	fmt.Println("All systems initialized. Waiting for tasks...")

	wg.Wait() // This will block until all tasks tracked by the WaitGroup have called Done()

	fmt.Println("All tasks completed or terminated.")
}
