package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type WalletResponse struct {
	// The resource name of the Balance.
	// Format: operations/{operation_id}
	Operation   string `json:"operation,omitempty"`
	DeviceGroup string `json:"deviceGroup,omitempty"`
}

type ListWalletsResponse struct {
	Wallets []*v1mpcwallets.MPCWallet `json:"wallets,omitempty"`
}

type ListAddressesResponse struct {
	Addresses []*v1mpcwallets.Address `json:"addresses,omitempty"`
}

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

func MpcWalletCreate(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("unable to read RegisterDevice request body: %v", err)
		httpGatewayTimeout(w)
		return
	}

	req := &v1mpcwallets.CreateMPCWalletRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal RegisterDevice request: %v", err)
		httpBadRequest(w)
		return
	}
	log.Debugf("creating wallet: %v", req)

	resp, err := mpcWalletServiceClient.CreateMPCWallet(r.Context(), req)
	if err != nil {
		log.Errorf("cannot create new wallet: %v", err)
		httpBadGateway(w)
		return
	}

	metadata, _ := resp.Metadata()

	wallet := &WalletResponse{
		Operation:   resp.Name(),
		DeviceGroup: metadata.GetDeviceGroup(),
	}

	if err := marhsallAndWriteJsonResponseWithOk(w, wallet); err != nil {
		log.Errorf("Cannot marshal and write create mpc wallet response: %v", err)
		httpBadGateway(w)
	}
}

func MpcWalletGenerateAddress(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("unable to read RegisterDevice request body: %v", err)
		httpGatewayTimeout(w)
		return
	}

	req := &v1mpcwallets.GenerateAddressRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal GenerateAddressRequest request: %v", err)
		httpBadRequest(w)
		return
	}
	log.Debugf("generating address: %v", req)

	resp, err := mpcWalletServiceClient.GenerateAddress(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot generate addres: %v", err)
		httpBadGateway(w)
		return
	}
	log.Debugf("generating address raw response: %v", resp)

	if err := marhsallAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("Cannot marshal and write generating address response: %v", err)
		httpBadGateway(w)
	}

}

func MpcWalletList(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Network id not passed to MpcWalletList")
		httpBadRequest(w)
		return
	}

	req := &v1mpcwallets.ListMPCWalletsRequest{
		Parent: fmt.Sprintf("pools/%s", poolId),
	}

	iter := mpcWalletServiceClient.ListMPCWallets(r.Context(), req)

	var wallets []*v1mpcwallets.MPCWallet
	for {
		wallet, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate wallets: %v", err)
			httpBadGateway(w)
			return
		}
		wallets = append(wallets, wallet)
	}

	response := &ListWalletsResponse{Wallets: wallets}

	if err := marhsallAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write mpc wallet list response: %v", err)
		httpBadGateway(w)
	}
}

func MpcAddressList(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("pool id not passed to MpcAddressList")
		httpBadRequest(w)
		return
	}

	networkId, found := vars["networkId"]
	if !found {
		log.Error("Network id not passed to MpcAddressList")
		httpBadRequest(w)
		return
	}

	mpcWalletId, found := vars["mpcWalletId"]
	if !found {
		log.Error("Network id not passed to MpcAddressList")
		httpBadRequest(w)
		return
	}

	req := &v1mpcwallets.ListAddressesRequest{
		Parent:    fmt.Sprintf("networks/%s", networkId),
		MpcWallet: fmt.Sprintf("pools/%s/mpcWallets/%s", poolId, mpcWalletId),
	}

	log.Debugf("calling listAddresses: %v", req)
	iter := mpcWalletServiceClient.ListAddresses(r.Context(), req)

	var addresses []*v1mpcwallets.Address
	for {
		address, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate addresses: %v", err)
			httpBadGateway(w)
			return
		}
		addresses = append(addresses, address)
	}

	log.Debugf("found addresses: %v", addresses)
	response := &ListAddressesResponse{Addresses: addresses}

	if err := marhsallAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write mpc address list response: %v", err)
		httpBadGateway(w)
	}
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

	if err := marhsallAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write mpc wallet list balances response: %v", err)
		httpBadGateway(w)
	}
}
