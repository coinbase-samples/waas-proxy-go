package mpc_key

import (
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
)

func RevokeDevice(w http.ResponseWriter, r *http.Request) {

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	req := &v1mpckeys.RevokeDeviceRequest{
		Name: string(body),
	}

	log.Debugf("revoke device request: %v", req)

	err = waas.GetClients().MpcKeyService.RevokeDevice(r.Context(), req)
	if err != nil {
		log.Errorf("cannot revoke device: %w", err)
		utils.HttpBadGateway(w)
		return
	}

	utils.HttpOk(w)
}
