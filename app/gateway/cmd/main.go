package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukkaal/GRPC-API/app/gateway/routes"
	"github.com/lukkaal/GRPC-API/app/gateway/rpc"
	"github.com/lukkaal/GRPC-API/config"
	"github.com/lukkaal/GRPC-API/pkg/utils/logger"
)

func main() {
	// init viper config
	config.InitConfig()

	// key: grpc.NewClient
	rpc.Init()

	// new gin gateway router
	router := routes.NewRouter()

	// set server
	server := &http.Server{
		Addr:           config.Conf.Server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// server init and listen
	go func() {
		fmt.Println("gateway listening on", config.Conf.Server.Port)
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.GinloggerObj.Fatalln("gateway start failed:", err)
		}
	}()

	// monitor signal
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit, os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM)
	sig := <-quit
	fmt.Println("exit signal received:", sig)

	// close server with timeout ctx
	ctx, cancel := context.WithTimeout(
		context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.GinloggerObj.Fatalln("HTTP server shutdown failed:", err)
	}

	fmt.Println("gateway exited gracefully")
}
