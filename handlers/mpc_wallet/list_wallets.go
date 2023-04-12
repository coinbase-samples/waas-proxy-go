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

type ListWalletsResponse struct {
	Wallets []*v1mpcwallets.MPCWallet `json:"wallets,omitempty"`
}

func ListWallets(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Network id not passed to MpcWalletList")
		utils.HttpBadRequest(w)
		return
	}

	req := &v1mpcwallets.ListMPCWalletsRequest{
		Parent: fmt.Sprintf("pools/%s", poolId),
	}

	iter := waas.GetClients().MpcWalletService.ListMPCWallets(r.Context(), req)

	var wallets []*v1mpcwallets.MPCWallet
	for {
		wallet, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate wallets: %v", err)
			utils.HttpBadGateway(w)
			return
		}
		wallets = append(wallets, wallet)
	}

	response := &ListWalletsResponse{Wallets: wallets}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write mpc wallet list response: %v", err)
		utils.HttpBadGateway(w)
	}
}