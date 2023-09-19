package server

import (
	"fmt"
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
		fmt.Println("\nGracefully shutting down...")
		srv.Shutdown(nil)
		os.Exit(0)
	}()
}
