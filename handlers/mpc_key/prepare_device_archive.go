package mpc_key

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

func PrepareDeviceArchive(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	deviceGroupId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceGroupId")
	if len(deviceGroupId) == 0 {
		return
	}

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	log.Debugf("raw PrepareDeviceArchive body - %v", string(body))
	req := &v1mpckeys.PrepareDeviceArchiveRequest{}
	if err := protojson.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal PrepareDeviceArchive request: %v", err)
		utils.HttpBadRequest(w)
		return
	}
	req.DeviceGroup = fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId)
	log.Debugf("preparing device archive: %v", req)

	resp, err := waas.GetClients().MpcKeyService.PrepareDeviceArchive(r.Context(), req)
	if err != nil {
		log.Errorf("cannot prepare device archive: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("prepare device archive raw response: %v", resp)

	meta, _ := resp.Metadata()
	log.Debugf("device archive metadata: %v", meta)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, meta); err != nil {
		log.Errorf("cannot marshal and write prepare device archive response: %v", err)
		utils.HttpBadGateway(w)
	}
}
