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
	// credentials
	router.HandleFunc("/v1/waas/proxy/credentials", RetrieveCredentials).Methods(http.MethodGet)

	// blockchain service
	router.HandleFunc("/v1/waas/proxy/networks", ListNetworks).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/networks/{networkId}/assets", ListAssets).Methods(http.MethodGet)

	// pool service
	router.HandleFunc("/v1/waas/proxy/pools", CreatePool).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/pools", ListPools).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/pools/{poolId}", GetPool).Methods(http.MethodGet)

	// mpc keys service
	router.HandleFunc("/v1/waas/proxy/mpckeys/registerDevice", MpcRegisterDevice).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}", MpcWalletListOperations).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/createSignature/{mpcKeyId}", MpcCreateSignature).Methods(http.MethodPost)

	// protocol service
	router.HandleFunc("/v1/waas/proxy/protocols/networks/{networkId}/tx/construct", ConstructTransaction).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/protocols/networks/{networkId}/tx/broadcast", BroadcastTransaction).Methods(http.MethodPost)

	//mpc wallets service
	router.HandleFunc("/v1/waas/proxy/mpcwallets", MpcWalletCreate).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/address", MpcWalletGenerateAddress).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/pools/{poolId}", MpcWalletList).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/networks/{networkId}/pools/{poolId}/mpcWallets/{mpcWalletId}/addresses", MpcAddressList).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/networks/{networkId}/addresses/{addressId}", MpcWalletListBalances).Methods(http.MethodGet)

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
