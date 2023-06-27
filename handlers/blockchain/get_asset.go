package blockchain

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	log "github.com/sirupsen/logrus"
)

func GetAsset(w http.ResponseWriter, r *http.Request) {

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
		return
	}

	assetId := utils.HttpPathVarOrSendBadRequest(w, r, "assetId")
	if len(assetId) == 0 {
		return
	}

	req := &v1blockchain.GetAssetRequest{
		Name: fmt.Sprintf("networks/%s/assets/%s", networkId, assetId),
	}

	log.Debugf("GetAsset request: %v", req)
	resp, err := waas.GetClients().BlockchainService.GetAsset(r.Context(), req)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	log.Debugf("GetAsset response: %v", resp)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("cannot marshal and write get device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}
