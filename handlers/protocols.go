package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/config"
	waasv1 "github.com/coinbase/waas-client-library-go/clients/v1"
	v1protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	v1types "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/types/v1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var protocolServiceClient *waasv1.ProtocolServiceClient

func initProtocolClient(ctx context.Context, config config.AppConfig) (err error) {

	opts := waasClientDefaults(config)

	if protocolServiceClient, err = waasv1.NewProtocolServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS protocol client: %w", err)
	}
	return
}

func BroadcastTransaction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	networkId, found := vars["networkId"]

	if !found {
		log.Error("Network id not passed to BroadcastTransaction")
		httpBadRequest(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read BroadcastTransaction request body: %v", err)
		httpGatewayTimeout(w)
		return
	}

	req := &v1protocols.BroadcastTransactionRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal BroadcastTransaction request: %v", err)
		httpBadRequest(w)
		return
	}

	req.Network = fmt.Sprintf("networks/%s", networkId)

	tx, err := protocolServiceClient.BroadcastTransaction(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot broadcast tx: %v", err)
		httpBadGateway(w)
		return
	}

	marshalTransactionAndWriteResponse(w, tx, http.StatusCreated)

}

func ConstructTransaction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	networkId, found := vars["networkId"]

	if !found {
		log.Error("Network id not passed to ConstructTransaction")
		httpBadRequest(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read ConstructTransaction request body: %v", err)
		httpGatewayTimeout(w)
		return
	}

	req := &v1protocols.ConstructTransactionRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal ConstructTransaction request: %v", err)
		httpBadRequest(w)
		return
	}

	req.Network = fmt.Sprintf("networks/%s", networkId)

	tx, err := protocolServiceClient.ConstructTransaction(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot construct transaction: %v", err)
		httpBadGateway(w)
		return
	}

	marshalTransactionAndWriteResponse(w, tx, http.StatusCreated)
}

func marshalTransactionAndWriteResponse(
	w http.ResponseWriter,
	tx *v1types.Transaction,
	status int,
) {
	body, err := json.Marshal(tx)
	if err != nil {
		log.Errorf("Cannot marshal tx struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatus(w, body, status); err != nil {
		log.Errorf("Cannot write tx response: %v", err)
		httpBadGateway(w)
		return
	}
}
