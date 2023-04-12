package mpc_key

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
)

func RegisterDevice(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read RegisterDevice request body: %v", err)
		utils.HttpGatewayTimeout(w)
		return
	}

	req := &v1mpckeys.RegisterDeviceRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal RegisterDevice request: %v", err)
		utils.HttpBadRequest(w)
		return
	}
	log.Debugf("registering device: %v", req)

	resp, err := waas.GetClients().MpcKeyService.RegisterDevice(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot register new device: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("register device raw response: %v", resp)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("Cannot marshal and write register device response: %v", err)
		utils.HttpBadGateway(w)
	}
}
