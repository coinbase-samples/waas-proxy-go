package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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

type BroadcastRequest struct {
	NetworkId         string `json:"network,omitempty"`
	SignedTransaction string `json:"signedTransaction,omitempty"`
}

func BroadcastTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read BroadcastTransaction request body: %v", err)
		httpGatewayTimeout(w)
		return
	}

	br := &BroadcastRequest{}
	if err := json.Unmarshal(body, br); err != nil {
		log.Errorf("Unable to unmarshal Broadcast request: %v", err)
		httpBadRequest(w)
		return
	}

	log.Infof("original signed: %s", br.SignedTransaction)
	encoded, _ := hex.DecodeString(br.SignedTransaction)

	log.Infof("encoded tx: %s", encoded)
	req := &v1protocols.BroadcastTransactionRequest{
		Network: br.NetworkId,
		Transaction: &v1types.Transaction{
			RawSignedTransaction: encoded,
		},
	}

	tx, err := protocolServiceClient.BroadcastTransaction(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot broadcast tx: %v", err)
		httpBadGateway(w)
		return
	}
	log.Debugf("broadcast result: %v", tx)

	if err := marhsallAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("Cannot marshal and wite broadcast tx response: %v", err)
		httpBadGateway(w)
	}
}

func ConstructTransaction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	networkId, found := vars["networkId"]

	if !found {
		log.Error("Network id not passed to ConstructTransaction")
		httpBadRequest(w)
		return
	}

	body, err := io.ReadAll(r.Body)
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

	if err := marhsallAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("Cannot marshal and wite construct tx response: %v", err)
		httpBadGateway(w)
	}
}
