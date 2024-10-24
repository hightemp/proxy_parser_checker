package main

import (
	"flag"

	"github.com/hightemp/proxy_parser_checker/internal/checker"
	"github.com/hightemp/proxy_parser_checker/internal/config"
	"github.com/hightemp/proxy_parser_checker/internal/logger"
	"github.com/hightemp/proxy_parser_checker/internal/parser"
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

	go parser.Loop(cfg)
	go checker.Loop(cfg)
}
