package main

import (
	"os"
	"os/signal"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/coinbase-samples/waas-proxy-go/proxy"
	log "github.com/sirupsen/logrus"
)

var (
	productIds = `["BTC-USD", "ETH-USD", "ADA-USD", "MATIC-USD", "ATOM-USD", "SOL-USD"]`
)

func main() {

	var app config.AppConfig

	if err := config.Setup(&app); err != nil {
		log.Fatalf("Unable to config app: %v", err)
	}

	config.LogInit(app)

	run := make(chan os.Signal, 1)
	signal.Notify(run, os.Interrupt)

	go proxy.ProcessMessages(app, run)

	<-run
}
