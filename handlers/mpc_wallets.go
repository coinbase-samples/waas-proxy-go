package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	waasv1 "github.com/coinbase/waas-client-library-go/clients/v1"

	"github.com/coinbase-samples/waas-proxy-go/config"
	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

var mpcWalletServiceClient *waasv1.MPCWalletServiceClient
var blockchainServiceClient *waasv1.BlockchainServiceClient

// Extension of API's balance
type Balance struct {
	// The resource name of the Balance.
	// Format: networks/{network_id}/addresses/{address_id}/balances/{balance_id}
	Name string `json:"name,omitempty"`
	// The resource name of the Asset to which this Balance corresponds.
	// Format: networks/{network}/assets/{asset}
	Asset string `json:"asset,omitempty"`
	// The amount of the Asset, denominated in atomic units of the asset (e.g., Wei for Ether),
	// as a base-10 number.
	Amount string `json:"amount,omitempty"`
	// The resource name of the MPCWallet to which this Balance belongs.
	// Format: pools/{pool}/mpcWallets/{mpcWallet}
	MpcWallet string `json:"mpc_wallet,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
	Decimals  int32  `json:"decimals,omitempty"`
}

type ListBalancesResponse struct {
	Balances []*Balance `json:"balances"`
}

func initMpcWalletClient(ctx context.Context, config config.AppConfig) (err error) {

	opts := waasClientDefaults(config)

	if mpcWalletServiceClient, err = waasv1.NewMPCWalletServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS mpc wallet client: %w", err)
	}
	if blockchainServiceClient, err = waasv1.NewBlockchainServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS blockchain client: %w", err)
	}
	return
}

func MpcWalletListBalances(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	vars := mux.Vars(r)

	networkId, found := vars["networkId"]
	if !found {
		log.Error("Network id not passed to MpcWalletListBalances")
		httpBadRequest(w)
		return
	}

	addressId, found := vars["addressId"]
	if !found {
		log.Error("Address id not passed to MpcWalletListBalances")
		httpBadRequest(w)
		return
	}

	req := &v1mpcwallets.ListBalancesRequest{
		Parent: fmt.Sprintf("networks/%s/addresses/%s", networkId, addressId),
	}

	iter := mpcWalletServiceClient.ListBalances(r.Context(), req)

	var balances []*v1mpcwallets.Balance
	for {
		balance, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate balances: %v", err)
			httpBadGateway(w)
			return
		}
		balances = append(balances, balance)
	}

	// TODO: switch to BatchGetAssets when ready
	var filledBalances []*Balance
	for i := 0; i < len(balances); i++ {
		b := balances[i]
		bReq := &v1blockchain.GetAssetRequest{
			Name: b.Asset,
		}
		asset, err := blockchainServiceClient.GetAsset(r.Context(), bReq)
		if err != nil {
			log.Warnf("Erroring fetching asset details - %v", err)
			continue
		}
		filledBalances = append(filledBalances, &Balance{
			Name:      b.Name,
			Asset:     b.Asset,
			Amount:    b.Amount,
			MpcWallet: b.MpcWallet,
			Symbol:    asset.AdvertisedSymbol,
			Decimals:  asset.Decimals,
		})
	}

	response := &ListBalancesResponse{Balances: filledBalances}

	body, err := json.Marshal(response)
	if err != nil {
		log.Errorf("Cannot marshal mpc wallet list balances struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatus(w, body, http.StatusOK); err != nil {
		log.Errorf("Cannot write list pools response: %v", err)
		httpBadGateway(w)
		return
	}
}
