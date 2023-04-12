package mpc_wallet

import (
	"encoding/json"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	log "github.com/sirupsen/logrus"
)

type WalletResponse struct {
	// The resource name of the Balance.
	// Format: operations/{operation_id}
	Operation   string `json:"operation,omitempty"`
	DeviceGroup string `json:"deviceGroup,omitempty"`
}

func CreateWallet(w http.ResponseWriter, r *http.Request) {

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	req := &v1mpcwallets.CreateMPCWalletRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal RegisterDevice request: %v", err)
		utils.HttpBadRequest(w)
		return
	}
	log.Debugf("creating wallet: %v", req)

	resp, err := waas.GetClients().MpcWalletService.CreateMPCWallet(r.Context(), req)
	if err != nil {
		log.Errorf("cannot create new wallet: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	metadata, _ := resp.Metadata()

	wallet := &WalletResponse{
		Operation:   resp.Name(),
		DeviceGroup: metadata.GetDeviceGroup(),
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, wallet); err != nil {
		log.Errorf("Cannot marshal and write create mpc wallet response: %v", err)
		utils.HttpBadGateway(w)
	}
}
