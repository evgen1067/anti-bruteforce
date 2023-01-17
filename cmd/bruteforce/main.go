package main

import (
	"flag"
	"github.com/evgen1067/anti-bruteforce/internal/app"
	"github.com/evgen1067/anti-bruteforce/internal/config"
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
	err = app.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
