package server

import (
	"github.com/ethereum/go-ethereum/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func HandleGracefulShutdown(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Trace("HandleGracefulShutdown", "server", "gracefully shutting down...")
		srv.Shutdown(nil)
		os.Exit(0)
	}()
}
