/**
 * Copyright 2023 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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

func PrepareDeviceBackup(w http.ResponseWriter, r *http.Request) {

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

	log.Debugf("raw PrepareDeviceBackup body - %v", string(body))
	req := &v1mpckeys.PrepareDeviceBackupRequest{}
	if err := protojson.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal PrepareDeviceBackup request: %v", err)
		utils.HttpBadRequest(w)
		return
	}
	req.DeviceGroup = fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId)
	log.Debugf("preparing device backup: %v", req)

	resp, err := waas.GetClients().MpcKeyService.PrepareDeviceBackup(r.Context(), req)
	if err != nil {
		log.Errorf("cannot prepare device backup: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("prepare device backup raw response: %v", resp)

	meta, _ := resp.Metadata()
	log.Debugf("device backup metadata: %v", meta)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, meta); err != nil {
		log.Errorf("cannot marshal and write prepare device backup response: %v", err)
		utils.HttpBadGateway(w)
	}
}
