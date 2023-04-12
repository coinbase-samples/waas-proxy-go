package pool

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

func CreatePool(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read CreatePool request body: %v", err)
		utils.HttpGatewayTimeout(w)
		return
	}

	req := &v1pools.CreatePoolRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal CreatePool request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	pool, err := waas.GetClients().PoolService.CreatePool(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot create pool: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, pool); err != nil {
		log.Errorf("Cannot marshal and write create pool response: %v", err)
		utils.HttpBadGateway(w)
	}
}
