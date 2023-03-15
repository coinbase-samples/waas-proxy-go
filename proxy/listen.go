package proxy

import (
	"io"
	"net/http"
	"os"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/coinbase-samples/waas-proxy-go/proxy/handlers"
	"github.com/gorilla/mux"
)

func ProcessMessages(config config.AppConfig, interrupt chan os.Signal) {

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "These aren't the droids you're looking for...\n")
	})

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ok\n")
	})

	router.HandleFunc("/v1/waas/proxy/credentials", handlers.RetrieveCredentials).Methods("GET")

	router.HandleFunc("/v1/waas/proxy/pool", handlers.CreatePool).Methods("PUT")

}
