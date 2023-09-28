package main

import (
	conf "bittoCralwer/config"
	"bittoCralwer/internal/common"
	"bittoCralwer/internal/model"
	"bittoCralwer/internal/server"
	"bittoCralwer/internal/service/blockchain/ether/blockscraper"
	"bittoCralwer/internal/service/blockchain/ether/txscraper"
	"flag"
	"github.com/ethereum/go-ethereum/log"
	"path"
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

	common.InitLog(common.Config{
		UseTerminal:        config.Log.Terminal.Use,
		UseFile:            config.Log.File.Use,
		TerminalJSONOutput: config.Log.Terminal.JSONFormat,
		VerbosityTerminal:  config.Log.Terminal.Verbosity,
		VerbosityFile:      config.Log.File.Verbosity,
		FilePath:           path.Join(config.Datadir.Root, config.Datadir.Log, config.Log.File.FileName),
	})

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		txscraper.StartScrapingTxs(config, repositories)
	}()

	srv := server.StartServer(config.Port.Http)

	// This function can also accept the waitgroup if you want to track the server goroutine
	server.HandleGracefulShutdown(srv)

	log.Trace("Initialized", "server", "all systems initialized. waiting for tasks...")

	wg.Wait() // This will block until all tasks tracked by the WaitGroup have called Done()

	log.Trace("Completed", "server", "all tasks completed or terminated.")
}
