package pool

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

func ListPools(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	req := &v1pools.ListPoolsRequest{}

	log.Debugf("ListPools request: %v", req)
	iter := waas.GetClients().PoolService.ListPools(r.Context(), req)

	var pools []*v1pools.Pool
	for {
		pool, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate pools: %v", err)
			utils.HttpBadGateway(w)
			return
		}

		pools = append(pools, pool)
	}

	response := &v1pools.ListPoolsResponse{Pools: pools}

	log.Debugf("ListPools response: %v", response)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write list pools response: %v", err)
		utils.HttpBadGateway(w)
	}
}
