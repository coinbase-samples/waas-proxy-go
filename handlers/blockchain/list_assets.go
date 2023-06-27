package blockchain

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
)

var pageSize = 50

func ListAssets(w http.ResponseWriter, r *http.Request) {

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
		return
	}

	names := r.URL.Query()["names"]
	if len(names) > 0 {
		BatchGetAsset(w, r, networkId, names)
		return
	}

	req := &v1blockchain.ListAssetsRequest{
		Parent: fmt.Sprintf("networks/%s", networkId),
	}

	filter := r.URL.Query().Get("filter")
	if len(filter) > 1 {
		req.Filter = filter
	}
	log.Debugf("listing assets: %v", req)

	iter := waas.GetClients().BlockchainService.ListAssets(r.Context(), req)

	var assets []*v1blockchain.Asset
	i := 1
	for i < pageSize {
		asset, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate assets: %v", err)
			utils.HttpBadGateway(w)
			return
		}

		assets = append(assets, asset)
		i++
	}

	response := &v1blockchain.ListAssetsResponse{Assets: assets}

	log.Debugf("listAssets response: %v", response)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write list assets response: %v", err)
		utils.HttpBadGateway(w)
	}
}
