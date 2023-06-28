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
)

type AddDeviceResponse struct {
	// The resource name of the Balance.
	// Format: operations/{operation_id}
	Operation            string   `json:"operation,omitempty"`
	DeviceGroup          string   `json:"deviceGroup,omitempty"`
	ParticipatingDevices []string `json:"participatingDevices,omitempty"`
}

func AddDevice(w http.ResponseWriter, r *http.Request) {

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

	req := &v1mpckeys.AddDeviceRequest{
		DeviceGroup: fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId),
		Device:      string(body),
	}

	log.Debugf("add device request: %v", req)

	resp, err := waas.GetClients().MpcKeyService.AddDevice(r.Context(), req)
	if err != nil {
		log.Errorf("cannot create device group: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	metadata, _ := resp.Metadata()

	response := &AddDeviceResponse{
		Operation:            resp.Name(),
		DeviceGroup:          metadata.GetDeviceGroup(),
		ParticipatingDevices: metadata.GetParticipatingDevices(),
	}

	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("cannot marshal and write create device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}
