package handlers

import (
	"io"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/coinbase-samples/waas-proxy-go/handlers/blockchain"
	"github.com/coinbase-samples/waas-proxy-go/handlers/credentials"
	"github.com/coinbase-samples/waas-proxy-go/handlers/mpc_key"
	"github.com/coinbase-samples/waas-proxy-go/handlers/mpc_transaction"
	"github.com/coinbase-samples/waas-proxy-go/handlers/mpc_wallet"
	"github.com/coinbase-samples/waas-proxy-go/handlers/pool"
	"github.com/coinbase-samples/waas-proxy-go/handlers/protocol"
	"github.com/gorilla/mux"
)

var appConfig config.AppConfig

func RegisterHandlers(config config.AppConfig, router *mux.Router) {

	appConfig = config

	registerDefaultHandlers(config, router)

	// credentials
	router.HandleFunc("/v1/waas/proxy/credentials", credentials.Retrieve).Methods(http.MethodGet)

	// Blockchain service
	router.HandleFunc("/v1/waas/proxy/networks", blockchain.ListNetworks).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/networks/{networkId}/assets", blockchain.ListAssets).Methods(http.MethodGet)

	// Pool service
	router.HandleFunc("/v1/waas/proxy/pools", pool.CreatePool).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/pools", pool.ListPools).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/pools/{poolId}", pool.GetPool).Methods(http.MethodGet)

	// MPC keys service
	router.HandleFunc("/v1/waas/proxy/mpckeys/registerDevice", mpc_key.RegisterDevice).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}", mpc_key.GetDeviceGroup).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/mpcOperations", mpc_key.ListOperations).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/createSignature/{mpcKeyId}", mpc_key.CreateSignature).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/prepareDeviceArchive", mpc_key.PrepareDeviceArchive).Methods(http.MethodPost)

	// Protocol service
	router.HandleFunc("/v1/waas/proxy/protocols/networks/{networkId}/tx/construct", protocol.ConstructTransaction).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/protocols/networks/{networkId}/tx/broadcast", protocol.BroadcastTransaction).Methods(http.MethodPost)

	// MPC wallets service
	router.HandleFunc("/v1/waas/proxy/mpcwallets", mpc_wallet.CreateWallet).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/address", mpc_wallet.GenerateAddress).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/pools/{poolId}", mpc_wallet.ListWallets).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/networks/{networkId}/pools/{poolId}/mpcWallets/{mpcWalletId}/addresses", mpc_wallet.ListAddresses).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/networks/{networkId}/addresses/{addressId}", mpc_wallet.ListBalances).Methods(http.MethodGet)

	// MPC Transactions
	router.HandleFunc("/v1/waas/proxy/mpctransactions/pools/{poolId}/mpcWallets/{mpcWalletId}", mpc_transaction.CreateMPCTransaction).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpctransactions/pools/{poolId}/mpcWallets/{mpcWalletId}", mpc_transaction.ListMpcTransactions).Methods(http.MethodGet)

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
