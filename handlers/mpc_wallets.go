package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	waasv1 "github.cbhq.net/cloud/waas-client-library-go/clients/v1"

	v1mpcwallets "github.cbhq.net/cloud/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

var mpcWalletServiceClient *waasv1.MPCWalletServiceClient

func initMpcWalletClient(ctx context.Context, config config.AppConfig) (err error) {

	opts := waasClientDefaults(config)

	if mpcWalletServiceClient, err = waasv1.NewMPCWalletServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS mpc wallet client: %w", err)
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

	response := &v1mpcwallets.ListBalancesResponse{Balances: balances}

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
