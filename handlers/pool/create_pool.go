package pool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

func CreatePool(w http.ResponseWriter, r *http.Request) {

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	req := &v1pools.CreatePoolRequest{}
	if err := protojson.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal CreatePool request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	pool, err := createPool(r.Context(), req)
	if err != nil {
		log.Error(err)
		utils.HttpBadGateway(w)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, pool); err != nil {
		log.Errorf("Cannot marshal and write create pool response: %v", err)
		utils.HttpBadGateway(w)
	}
}

func createPool(
	ctx context.Context,
	req *v1pools.CreatePoolRequest,
) (*v1pools.Pool, error) {
	pool, err := waas.GetClients().PoolService.CreatePool(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Cannot create pool: %w", err)
	}
	return pool, nil
}
