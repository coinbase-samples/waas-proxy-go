package pool

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

func GetPool(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	poolId, found := vars["poolId"]

	if !found {
		log.Error("Pool id not passed to GetPool")
		utils.HttpBadRequest(w)
		return
	}

	req := &v1pools.GetPoolRequest{
		Name: fmt.Sprintf("pools/%s", poolId),
	}

	pool, err := waas.GetClients().PoolService.GetPool(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot get pool: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, pool); err != nil {
		log.Errorf("Cannot marshal and write get pool response: %v", err)
		utils.HttpBadGateway(w)
	}
}
