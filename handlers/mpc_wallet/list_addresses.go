package mpc_wallet

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

func ListAddresses(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
		return
	}

	mpcWalletId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcWalletId")
	if len(mpcWalletId) == 0 {
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
	response := &v1mpcwallets.ListAddressesResponse{Addresses: addresses}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write mpc address list response: %v", err)
		utils.HttpBadGateway(w)
	}
}
