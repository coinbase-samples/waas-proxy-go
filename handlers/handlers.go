/**
 * Copyright 2023 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package handlers

import (
	"io"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/coinbase-samples/waas-proxy-go/handlers/backup"
	"github.com/coinbase-samples/waas-proxy-go/handlers/blockchain"
	"github.com/coinbase-samples/waas-proxy-go/handlers/combined"
	"github.com/coinbase-samples/waas-proxy-go/handlers/mpc_key"
	"github.com/coinbase-samples/waas-proxy-go/handlers/mpc_transaction"
	"github.com/coinbase-samples/waas-proxy-go/handlers/mpc_wallet"
	"github.com/coinbase-samples/waas-proxy-go/handlers/pool"
	"github.com/coinbase-samples/waas-proxy-go/handlers/protocol"
	"github.com/gorilla/mux"
)

func RegisterHandlers(config config.AppConfig, router *mux.Router) {

	registerDefaultHandlers(config, router)

	// Blockchain service
	router.HandleFunc("/v1/waas/proxy/networks", blockchain.ListNetworks).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/networks/{networkId}/assets", blockchain.ListAssets).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/networks/{networkId}/assets/{assetId}", blockchain.GetAsset).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/networks/{networkId}", blockchain.GetNetwork).Methods(http.MethodGet)

	// Pool service
	router.HandleFunc("/v1/waas/proxy/pools", pool.CreatePool).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/pools", pool.ListPools).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/pools/{poolId}", pool.GetPool).Methods(http.MethodGet)

	// MPC keys service
	router.HandleFunc("/v1/waas/proxy/mpckeys/devices/{deviceId}", mpc_key.GetDevice).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/registerDevice", mpc_key.RegisterDevice).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}", mpc_key.GetDeviceGroup).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/mpcOperations", mpc_key.ListOperations).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/mpcKeys/{mpcKeyId}/createSignature", mpc_key.CreateSignature).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/prepareDeviceArchive", mpc_key.PrepareDeviceArchive).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/prepareDeviceBackup", mpc_key.PrepareDeviceBackup).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/addDevice", mpc_key.AddDevice).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/device/revoke", mpc_key.RevokeDevice).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpckeys/pools/{poolId}/deviceGroups/{deviceGroupId}/mpcKeys/{mpcKeyId}", mpc_key.GetMpcKey).Methods(http.MethodGet)

	// Protocol service
	router.HandleFunc("/v1/waas/proxy/protocols/networks/{networkId}/tx/construct", protocol.ConstructTransaction).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/protocols/networks/{networkId}/tx/constructTransfer", protocol.ConstructTransferTransaction).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/protocols/networks/{networkId}/tx/broadcast", protocol.BroadcastTransaction).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/protocols/networks/{networkId}/estimateFee", protocol.EstimateFee).Methods(http.MethodPost)

	// MPC wallets service
	router.HandleFunc("/v1/waas/proxy/mpcwallets/pools/{poolId}/mpcWallets/{mpcWalletId}", mpc_wallet.GetWallet).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpcwallets", mpc_wallet.CreateWallet).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/wait", mpc_wallet.WaitWallet).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/pools/{poolId}", mpc_wallet.ListWallets).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/address", mpc_wallet.GenerateAddress).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/networks/{networkId}/addresses/{addressId}", mpc_wallet.GetAddress).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/networks/{networkId}/pools/{poolId}/mpcWallets/{mpcWalletId}/addresses", mpc_wallet.ListAddresses).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpcwallets/networks/{networkId}/addresses/{addressId}/balances", mpc_wallet.ListBalances).Methods(http.MethodGet)

	// MPC Transactions
	router.HandleFunc("/v1/waas/proxy/mpctransactions/pools/{poolId}/mpcWallets/{mpcWalletId}", mpc_transaction.CreateMPCTransaction).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/mpctransactions/pools/{poolId}/mpcWallets/{mpcWalletId}", mpc_transaction.ListMpcTransactions).Methods(http.MethodGet)
	router.HandleFunc("/v1/waas/proxy/mpctransactions/pools/{poolId}/mpcWallets/{mpcWalletId}/mpcTransactions/{mpcTransactionId}", mpc_transaction.GetMpcTransaction).Methods(http.MethodGet)

	// Combined
	router.HandleFunc("/v1/waas/proxy/combined/pools/{poolId}/deviceGroups/{deviceGroupId}/mpcKeys/{mpcKeyId}/constructAndSign", combined.ConstructAndSign).Methods(http.MethodPost)
	router.HandleFunc("/v1/waas/proxy/combined/waitAndBroadcast", combined.WaitSignAndBroadcast).Methods(http.MethodPost)

	if config.S3Enabled {
		// S3 Urls
		router.HandleFunc("/v1/waas/proxy/backup/method/{s3Method}/objectKey/{objectKey}", backup.GetPresignedS3Url).Methods(http.MethodGet)
	}
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
