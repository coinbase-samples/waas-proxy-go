package mpc_key

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
)

func GetDevice(w http.ResponseWriter, r *http.Request) {

	deviceId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceId")
	if len(deviceId) == 0 {
		return
	}

	req := &v1mpckeys.GetDeviceRequest{
		Name: fmt.Sprintf("devices/%s", deviceId),
	}

	log.Debugf("GetDevice request: %v", req)
	resp, err := waas.GetClients().MpcKeyService.GetDevice(r.Context(), req)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	log.Debugf("GetDevice response: %v", resp)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("Cannot marshal and write get device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}
