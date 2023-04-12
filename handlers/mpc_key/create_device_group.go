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

type CreateDeviceGroupResponse struct {
	// The resource name of the Balance.
	// Format: operations/{operation_id}
	Operation   string `json:"operation,omitempty"`
	DeviceGroup string `json:"deviceGroup,omitempty"`
}

func CreateDeviceGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Pool id not passed to MpcWalletCreate")
		utils.HttpBadRequest(w)
		return
	}

	deviceId, found := vars["deviceId"]
	if !found {
		log.Error("Device Id not passed to MpcWalletListOperations")
		utils.HttpBadRequest(w)
		return
	}

	req := &v1mpckeys.CreateDeviceGroupRequest{
		Parent: fmt.Sprintf("pools/%s/device/%s", poolId, deviceId),
	}

	resp, err := waas.GetClients().MpcKeyService.CreateDeviceGroup(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot list mpc operations: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	metadata, _ := resp.Metadata()

	finalResp := &CreateDeviceGroupResponse{
		Operation:   resp.Name(),
		DeviceGroup: metadata.GetDeviceGroup(),
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, finalResp); err != nil {
		log.Errorf("Cannot marshal and write create device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}
