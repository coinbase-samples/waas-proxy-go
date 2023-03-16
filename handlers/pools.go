package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	poolspb "github.com/WaaS-Private-Preview-v1/waas-client-library/go/coinbase/cloud/pools/v1alpha1"

	"github.com/WaaS-Private-Preview-v1/waas-client-library/go/coinbase/cloud/clients"
)

var poolServiceClient *clients.PoolServiceClient

func initPoolClient(ctx context.Context, config config.AppConfig) (err error) {

	endpoint, opts := waasClientDefaults(config, "waas/pools")

	if poolServiceClient, err = clients.NewV1Alpha1PoolServiceClient(
		ctx,
		endpoint,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS pool client: %w", err)
	}
	return
}

func GetPool(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	poolId, found := vars["poolId"]

	if !found {
		log.Error("Pool id not passed to GetPool")
		httpBadRequest(w)
		return
	}

	req := &poolspb.GetPoolRequest{
		Name: fmt.Sprintf("pools/%s", poolId),
	}

	pool, err := poolServiceClient.GetPool(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot get pool: %v", err)
		httpBadGateway(w)
		return
	}

	marshalPoolAndWriteResponse(w, pool, http.StatusOK)
}

func CreatePool(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read CreatePool request body: %v", err)
		httpGatewayTimeout(w)
		return
	}

	req := &poolspb.CreatePoolRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal CreatePool request: %v", err)
		httpBadRequest(w)
		return
	}

	pool, err := poolServiceClient.CreatePool(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot create pool: %v", err)
		httpBadGateway(w)
		return
	}

	marshalPoolAndWriteResponse(w, pool, http.StatusCreated)
}

func marshalPoolAndWriteResponse(
	w http.ResponseWriter,
	pool *poolspb.Pool,
	status int,
) {
	body, err := json.Marshal(pool)
	if err != nil {
		log.Errorf("Cannot marshal pool struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatus(w, body, status); err != nil {
		log.Errorf("Cannot write pool response: %v", err)
		httpBadGateway(w)
		return
	}
}
