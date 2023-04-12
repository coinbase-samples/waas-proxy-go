package mpc_key

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
)

type CreateDeviceGroupResponse struct {
	// The resource name of the Balance.
	// Format: operations/{operation_id}
	Operation   string `json:"operation,omitempty"`
	DeviceGroup string `json:"deviceGroup,omitempty"`
}

func CreateDeviceGroup(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	deviceId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceId")
	if len(deviceId) == 0 {
		return
	}

	response, err := createDeviceGroup(r.Context(), poolId, deviceId)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write create device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}

func createDeviceGroup(
	ctx context.Context,
	poolId,
	deviceId string,
) (*CreateDeviceGroupResponse, error) {

	req := &v1mpckeys.CreateDeviceGroupRequest{
		Parent: fmt.Sprintf("pools/%s/device/%s", poolId, deviceId),
	}

	resp, err := waas.GetClients().MpcKeyService.CreateDeviceGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Cannot list mpc operations: %", err)
	}

	metadata, _ := resp.Metadata()

	response := &CreateDeviceGroupResponse{
		Operation:   resp.Name(),
		DeviceGroup: metadata.GetDeviceGroup(),
	}

	return response, nil
}
