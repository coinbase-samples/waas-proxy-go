package blockchain

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
)

func ListAssets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	networkId, found := vars["networkId"]
	if !found {
		log.Error("networkId not passed to ListAssets")
		utils.HttpBadRequest(w)
		return
	}

	// TODO: This needs to page for the end client - iterator blasts through everything

	req := &v1blockchain.ListAssetsRequest{
		Parent: fmt.Sprintf("networks/%s", networkId),
	}

	filter := r.URL.Query().Get("filter")
	if len(filter) > 1 {
		req.Filter = filter
	}
	log.Infof("requesting assets: %v", req)

	iter := waas.GetClients().BlockchainService.ListAssets(r.Context(), req)

	var assets []*v1blockchain.Asset
	i := 1
	for i < 50 {
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

	log.Debugf("raw listAssets response: %s", response.String())
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write list assets response: %v", err)
		utils.HttpBadGateway(w)
	}
}
