package mpc_key

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func ListOperations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Pool id not passed to ListOperations")
		utils.HttpBadRequest(w)
		return
	}

	deviceGroupId, found := vars["deviceGroupId"]
	if !found {
		log.Error("Device Group Id not passed to ListOperations")
		utils.HttpBadRequest(w)
		return
	}

	req := &v1mpckeys.ListMPCOperationsRequest{
		Parent: fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId),
	}

	log.Debugf("listing mpc op request: %v", req)
	resp, err := waas.GetClients().MpcKeyService.ListMPCOperations(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot list mpc operations: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("Cannot marshal and write mpc operations response: %v", err)
		utils.HttpBadGateway(w)
	}
}
