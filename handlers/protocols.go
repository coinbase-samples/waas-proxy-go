package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	waasv1 "github.cbhq.net/cloud/waas-client-library-go/clients/v1"
	v1protocols "github.cbhq.net/cloud/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"

	"github.com/coinbase-samples/waas-proxy-go/config"
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
		log.Errorf("Cannot create pool: %v", err)
		httpBadGateway(w)
		return
	}

	body, err = json.Marshal(tx)
	if err != nil {
		log.Errorf("Cannot marshal transaction struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatus(w, body, http.StatusCreated); err != nil {
		log.Errorf("Cannot write create transaction response: %v", err)
		httpBadGateway(w)
		return
	}

}
