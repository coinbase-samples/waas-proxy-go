package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/config"
	waasv1 "github.com/coinbase/waas-client-library-go/clients/v1"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var mpcKeysServiceClient *waasv1.MPCKeyServiceClient

func initMpcKeyClient(ctx context.Context, config config.AppConfig) (err error) {

	opts := waasClientDefaults(config)

	if mpcKeysServiceClient, err = waasv1.NewMPCKeyServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS mpc key client: %w", err)
	}

	return
}

func MpcWalletListOperations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Pool id not passed to MpcWalletListOperations")
		httpBadRequest(w)
		return
	}

	deviceGroupId, found := vars["deviceGroupId"]
	if !found {
		log.Error("Device Group Id not passed to MpcWalletListOperations")
		httpBadRequest(w)
		return
	}

	req := &v1mpckeys.ListMPCOperationsRequest{
		Parent: fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId),
	}

	resp, err := mpcKeysServiceClient.ListMPCOperations(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot list mpc operations: %v", err)
		httpBadGateway(w)
		return
	}
	body, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Cannot marshal tx struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatus(w, body, http.StatusOK); err != nil {
		log.Errorf("Cannot write tx response: %v", err)
		httpBadGateway(w)
		return
	}

}
