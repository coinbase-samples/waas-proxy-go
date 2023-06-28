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
	"context"
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
)

type CreateDeviceGroupResponse struct {
	// The resource name of the Balance.
	// Format: operations/{operation_id}
	Operation   string `json:"operation,omitempty"`
	DeviceGroup string `json:"deviceGroup,omitempty"`
}

func CreateDeviceGroup(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	deviceId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceId")
	if len(deviceId) == 0 {
		return
	}

	response, err := createDeviceGroup(r.Context(), poolId, deviceId)
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

func createDeviceGroup(
	ctx context.Context,
	poolId,
	deviceId string,
) (*CreateDeviceGroupResponse, error) {

	req := &v1mpckeys.CreateDeviceGroupRequest{
		Parent: fmt.Sprintf("pools/%s/device/%s", poolId, deviceId),
	}

	log.Debugf("create device group request: %v", req)
	resp, err := waas.GetClients().MpcKeyService.CreateDeviceGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot create device group: %w", err)
	}

	metadata, _ := resp.Metadata()

	response := &CreateDeviceGroupResponse{
		Operation:   resp.Name(),
		DeviceGroup: metadata.GetDeviceGroup(),
	}
	log.Debugf("create device group response: %v", response)
	return response, nil
}
