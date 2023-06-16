package blockchain

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	log "github.com/sirupsen/logrus"
)

func GetNetwork(w http.ResponseWriter, r *http.Request) {

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
		return
	}

	req := &v1blockchain.GetNetworkRequest{
		Name: fmt.Sprintf("networks/%s", networkId),
	}

	resp, err := waas.GetClients().BlockchainService.GetNetwork(r.Context(), req)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("Cannot marshal and write get device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}