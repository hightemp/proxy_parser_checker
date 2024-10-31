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

const (
	VERSION = "v0.0.1"
)

func init() {
	logger.InitLogger()
}

func main() {
	logger.LogInfo("proxy_parser_checker Version: %s", VERSION)

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
