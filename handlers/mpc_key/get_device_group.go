package mpc_key

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
)

func GetDeviceGroup(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	deviceGroupId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceGroupId")
	if len(deviceGroupId) == 0 {
		return
	}

	req := &v1mpckeys.GetDeviceGroupRequest{
		Name: fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId),
	}

	log.Debugf("GetDeviceGroup request: %v", req)
	resp, err := waas.GetClients().MpcKeyService.GetDeviceGroup(r.Context(), req)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	log.Debugf("GetDeviceGroup response: %v", resp)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("cannot marshal and write get device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}
