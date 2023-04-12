package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

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

	if err := marhsallAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write list pools response: %v", err)
		httpBadGateway(w)
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

	if err := marhsallAndWriteJsonResponseWithOk(w, pool); err != nil {
		log.Errorf("Cannot marshal and write get pool response: %v", err)
		httpBadGateway(w)
	}
}

func CreatePool(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
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

	if err := marhsallAndWriteJsonResponseWithOk(w, pool); err != nil {
		log.Errorf("Cannot marshal and write create pool response: %v", err)
		httpBadGateway(w)
	}
}
