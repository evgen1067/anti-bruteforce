package main

import (
	"flag"
	"github.com/evgen1067/anti-bruteforce/internal/app"
	"github.com/evgen1067/anti-bruteforce/internal/config"
	"github.com/evgen1067/anti-bruteforce/internal/logger"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/local.json", "Path to json configuration file")
}

func main() {
	flag.Parse()
	cfg, err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}
	zLog, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = app.Run(zLog, cfg)
	if err != nil {
		zLog.Error(err.Error())
	}
}
