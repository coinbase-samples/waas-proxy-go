package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
)

func ListNetworks(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	req := &v1blockchain.ListNetworksRequest{}

	iter := blockchainServiceClient.ListNetworks(r.Context(), req)

	var networks []*v1blockchain.Network
	for {
		network, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate networks: %v", err)
			httpBadGateway(w)
			return
		}

		networks = append(networks, network)
	}

	response := &v1blockchain.ListNetworksResponse{Networks: networks}

	if err := marhsallAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write list networks response: %v", err)
		httpBadGateway(w)
	}
}

func ListAssets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	networkId, found := vars["networkId"]
	if !found {
		log.Error("networkId not passed to ListAssets")
		httpBadRequest(w)
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

	iter := blockchainServiceClient.ListAssets(r.Context(), req)

	var assets []*v1blockchain.Asset
	i := 1
	for i < 50 {
		asset, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate assets: %v", err)
			httpBadGateway(w)
			return
		}

		assets = append(assets, asset)
		i++
	}

	response := &v1blockchain.ListAssetsResponse{Assets: assets}

	log.Debugf("raw listAssets response: %s", response.String())
	if err := marhsallAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write list assets response: %v", err)
		httpBadGateway(w)
	}
}
