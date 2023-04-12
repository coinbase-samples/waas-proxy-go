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

	// TODO: This needs to page for the end client - iterator blasts through everything

	req := &v1blockchain.ListNetworksRequest{}

	iter := waas.GetClients().BlockchainService.ListNetworks(r.Context(), req)

	var networks []*v1blockchain.Network
	for {
		network, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate networks: %v", err)
			utils.HttpBadGateway(w)
			return
		}

		networks = append(networks, network)
	}

	response := &v1blockchain.ListNetworksResponse{Networks: networks}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write list networks response: %v", err)
		utils.HttpBadGateway(w)
	}
}