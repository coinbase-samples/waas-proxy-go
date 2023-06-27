package blockchain

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	log "github.com/sirupsen/logrus"
)

func BatchGetAsset(w http.ResponseWriter, r *http.Request, networkId string, names []string) {

	req := &v1blockchain.BatchGetAssetsRequest{
		Parent: fmt.Sprintf("networks/%s", networkId),
		Names:  names,
	}

	log.Debugf("BatchGetAssets request: %v", req)
	resp, err := waas.GetClients().BlockchainService.BatchGetAssets(r.Context(), req)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	log.Debugf("BatchGetAssets response: %v", resp)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("cannot marshal and write get device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}
