package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/hightemp/proxy_parser_checker/internal/checker"
	"github.com/hightemp/proxy_parser_checker/internal/config"
	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"github.com/hightemp/proxy_parser_checker/internal/parser"
	"github.com/hightemp/proxy_parser_checker/internal/server"
)

func main() {
	logger.InitLogger()

	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	err := config.Load(*configPath)

	if err != nil {
		logger.PanicError("%v", err)
	}

	cfg := config.GetConfig()
	logger.LogDebug("Config loaded")

	go server.Start()

	go parser.Loop(cfg)
	go checker.Loop(cfg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan

	logger.LogDebug("Received signal: %v", sig)
}
