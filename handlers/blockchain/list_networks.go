package blockchain

import (
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
)

func ListNetworks(w http.ResponseWriter, r *http.Request) {

	req := &v1blockchain.ListNetworksRequest{}

	log.Debugf("ListNetworks request: %v", req)
	requestPageInfo, err := utils.HttpRequestPageInfo(r)
	if err != nil {
		utils.HttpBadRequest(w)
		return
	}

	if requestPageInfo.Passed() {
		req.PageToken = requestPageInfo.Token
		req.PageSize = requestPageInfo.Size
	}

	iter := waas.GetClients().BlockchainService.ListNetworks(r.Context(), req)

	_, err = iter.Next()
	if err != nil && err != iterator.Done {
		log.Errorf("cannot iterate assets: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	response := iter.Response()

	log.Debugf("ListNetworks response: %v", response)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("cannot marshal and write list networks response: %v", err)
		utils.HttpBadGateway(w)
	}
}
