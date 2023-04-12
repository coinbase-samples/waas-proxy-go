package mpc_wallet

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type ListAddressesResponse struct {
	Addresses []*v1mpcwallets.Address `json:"addresses,omitempty"`
}

func ListAddresses(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("pool id not passed to MpcAddressList")
		utils.HttpBadRequest(w)
		return
	}

	networkId, found := vars["networkId"]
	if !found {
		log.Error("Network id not passed to MpcAddressList")
		utils.HttpBadRequest(w)
		return
	}

	mpcWalletId, found := vars["mpcWalletId"]
	if !found {
		log.Error("Network id not passed to MpcAddressList")
		utils.HttpBadRequest(w)
		return
	}

	req := &v1mpcwallets.ListAddressesRequest{
		Parent:    fmt.Sprintf("networks/%s", networkId),
		MpcWallet: fmt.Sprintf("pools/%s/mpcWallets/%s", poolId, mpcWalletId),
	}

	log.Debugf("calling listAddresses: %v", req)
	iter := waas.GetClients().MpcWalletService.ListAddresses(r.Context(), req)

	var addresses []*v1mpcwallets.Address
	for {
		address, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate addresses: %v", err)
			utils.HttpBadGateway(w)
			return
		}
		addresses = append(addresses, address)
	}

	log.Debugf("found addresses: %v", addresses)
	response := &ListAddressesResponse{Addresses: addresses}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write mpc address list response: %v", err)
		utils.HttpBadGateway(w)
	}
}
