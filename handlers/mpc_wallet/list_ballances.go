package mpc_wallet

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

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
	MpcWallet  string                        `json:"mpc_wallet,omitempty"`
	Symbol     string                        `json:"symbol,omitempty"`
	Decimals   int32                         `json:"decimals,omitempty"`
	Definition v1blockchain.Asset_Definition `json:"definition,omitempty"`
}

type ListBalancesResponse struct {
	Balances []*Balance `json:"balances"`
}

func ListBalances(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
		return
	}

	addressId := utils.HttpPathVarOrSendBadRequest(w, r, "addressId")
	if len(addressId) == 0 {
		return
	}

	req := &v1mpcwallets.ListBalancesRequest{
		Parent: fmt.Sprintf("networks/%s/addresses/%s", networkId, addressId),
	}

	iter := waas.GetClients().MpcWalletService.ListBalances(r.Context(), req)

	var balances []*v1mpcwallets.Balance
	for {
		balance, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate balances: %v", err)
			utils.HttpBadGateway(w)
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
		asset, err := waas.GetClients().BlockchainService.GetAsset(r.Context(), bReq)
		if err != nil {
			log.Warnf("Erroring fetching asset details - %v", err)
			continue
		}
		filledBalances = append(filledBalances, &Balance{
			Name:       b.Name,
			Asset:      b.Asset,
			Amount:     b.Amount,
			MpcWallet:  b.MpcWallet,
			Symbol:     asset.AdvertisedSymbol,
			Decimals:   asset.Decimals,
			Definition: *asset.Definition,
		})
	}

	response := &ListBalancesResponse{Balances: filledBalances}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write mpc wallet list balances response: %v", err)
		utils.HttpBadGateway(w)
	}
}
