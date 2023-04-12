package mpc_wallet

import (
	"encoding/json"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	log "github.com/sirupsen/logrus"
)

func GenerateAddress(w http.ResponseWriter, r *http.Request) {

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	req := &v1mpcwallets.GenerateAddressRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal GenerateAddressRequest request: %v", err)
		utils.HttpBadRequest(w)
		return
	}
	log.Debugf("generating address: %v", req)

	resp, err := waas.GetClients().MpcWalletService.GenerateAddress(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot generate addres: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("generating address raw response: %v", resp)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("Cannot marshal and write generating address response: %v", err)
		utils.HttpBadGateway(w)
	}

}
