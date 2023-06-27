package mpc_key

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
)

func GetMpcKey(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	deviceGroupId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceGroupId")
	if len(deviceGroupId) == 0 {
		return
	}

	mpcKeyId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcKeyId")
	if len(mpcKeyId) == 0 {
		return
	}

	req := &v1mpckeys.GetMPCKeyRequest{
		Name: fmt.Sprintf("pools/%s/deviceGroups/%s/mpcKeys/%s", poolId, deviceGroupId, mpcKeyId),
	}

	resp, err := waas.GetClients().MpcKeyService.GetMPCKey(r.Context(), req)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("cannot marshal and write get device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}
