package pool

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

func GetPool(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	req := &v1pools.GetPoolRequest{
		Name: fmt.Sprintf("pools/%s", poolId),
	}

	log.Debugf("GetPool request: %v", req)
	pool, err := waas.GetClients().PoolService.GetPool(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot get pool: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("GetPool response: %v", pool)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, pool); err != nil {
		log.Errorf("Cannot marshal and write get pool response: %v", err)
		utils.HttpBadGateway(w)
	}
}
