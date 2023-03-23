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
	"google.golang.org/api/iterator"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"

	waasv1 "github.com/coinbase/waas-client-library-go/clients/v1"
)

var poolServiceClient *waasv1.PoolServiceClient

func initPoolClient(ctx context.Context, config config.AppConfig) (err error) {

	opts := waasClientDefaults(config)

	if poolServiceClient, err = waasv1.NewPoolServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS pool client: %w", err)
	}
	return
}

func ListPools(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	req := &v1pools.ListPoolsRequest{}

	iter := poolServiceClient.ListPools(r.Context(), req)

	var pools []*v1pools.Pool
	for {
		pool, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate pools: %v", err)
			httpBadGateway(w)
			return
		}

		pools = append(pools, pool)
	}

	response := &v1pools.ListPoolsResponse{Pools: pools}

	body, err := json.Marshal(response)
	if err != nil {
		log.Errorf("Cannot marshal list pools struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatus(w, body, http.StatusOK); err != nil {
		log.Errorf("Cannot write list pools response: %v", err)
		httpBadGateway(w)
		return
	}
}

func GetPool(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	poolId, found := vars["poolId"]

	if !found {
		log.Error("Pool id not passed to GetPool")
		httpBadRequest(w)
		return
	}

	req := &v1pools.GetPoolRequest{
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

	req := &v1pools.CreatePoolRequest{}
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
	pool *v1pools.Pool,
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
