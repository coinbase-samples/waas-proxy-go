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
	"time"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	log "github.com/sirupsen/logrus"
)

func ListOperations(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	deviceGroupId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceGroupId")
	if len(deviceGroupId) == 0 {
		return
	}

	response, err := listOperations(r.Context(), poolId, deviceGroupId)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("cannot marshal and write mpc operations response: %v", err)
		utils.HttpBadGateway(w)
	}
}

func listOperations(
	ctx context.Context,
	poolId,
	deviceGroupId string,
) (*v1mpckeys.ListMPCOperationsResponse, error) {

	req := &v1mpckeys.ListMPCOperationsRequest{
		Parent: fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId),
	}

	log.Debugf("listing mpc op request: %v", req)

	response, err := waas.GetClients().MpcKeyService.ListMPCOperations(ctx, req)
	if err != nil || len(response.MpcOperations) < 1 {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 200)
			response, err = waas.GetClients().MpcKeyService.ListMPCOperations(ctx, req)
			log.Debugf("fetched operations: %v", response)
			if err == nil && len(response.MpcOperations) > 0 {
				return response, nil
			}
		}
		if err != nil {
			return nil, fmt.Errorf("cannot list mpc operations: %w", err)
		}
	}

	return response, nil
}
