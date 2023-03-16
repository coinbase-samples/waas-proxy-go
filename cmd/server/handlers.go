package main

import (
	"io"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/gorilla/mux"

	"github.com/coinbase-samples/waas-proxy-go/handlers"
)

func registerHandlers(config config.AppConfig, router *mux.Router) {

	registerDefaultHandlers(config, router)

	router.HandleFunc("/v1/waas/proxy/credentials", handlers.RetrieveCredentials).Methods("GET")

	router.HandleFunc("/v1/waas/proxy/pool", handlers.CreatePool).Methods("PUT")

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
