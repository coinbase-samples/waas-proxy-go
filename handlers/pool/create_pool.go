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
package pool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

func CreatePool(w http.ResponseWriter, r *http.Request) {

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	req := &v1pools.CreatePoolRequest{}
	if err := protojson.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal CreatePool request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	log.Debugf("CreatePool request: %v", req)
	pool, err := createPool(r.Context(), req)
	if err != nil {
		log.Error(err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("CreatePool response: %v", pool)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, pool); err != nil {
		log.Errorf("cannot marshal and write create pool response: %v", err)
		utils.HttpBadGateway(w)
	}
}

func createPool(
	ctx context.Context,
	req *v1pools.CreatePoolRequest,
) (*v1pools.Pool, error) {
	pool, err := waas.GetClients().PoolService.CreatePool(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cannot create pool: %w", err)
	}
	return pool, nil
}
