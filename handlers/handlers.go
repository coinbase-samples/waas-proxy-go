package handlers

import (
	"io"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var appConfig config.AppConfig

func RegisterHandlers(config config.AppConfig, router *mux.Router) {

	appConfig = config

	if err := initWaaSClients(config); err != nil {
		log.Fatalf("Unable to init WaaS clients: %v", err)
	}

	registerDefaultHandlers(config, router)

	router.HandleFunc("/v1/waas/proxy/credentials", RetrieveCredentials).Methods(http.MethodGet)

	router.HandleFunc("/v1/waas/proxy/pools", CreatePool).Methods(http.MethodPut)

	router.HandleFunc("/v1/waas/proxy/pools", ListPools).Methods(http.MethodGet)

	//router.HandleFunc("/v1/waas/proxy/pools/{poolId}/wallets/{walletId}", GetWallet).Methods(http.MethodGet)

	router.HandleFunc("/v1/waas/proxy/pools/{poolId}", GetPool).Methods(http.MethodGet)

}

func registerDefaultHandlers(config config.AppConfig, router *mux.Router) {

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "These aren't the droids you're looking for...\n")
	})

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ok\n")
	})

}
