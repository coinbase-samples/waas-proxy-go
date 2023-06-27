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
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"

	v1pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
)

func GetPool(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	req := &v1pools.GetPoolRequest{
		Name: fmt.Sprintf("pools/%s", poolId),
	}

	log.Debugf("GetPool request: %v", req)
	pool, err := waas.GetClients().PoolService.GetPool(r.Context(), req)
	if err != nil {
		log.Errorf("cannot get pool: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("GetPool response: %v", pool)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, pool); err != nil {
		log.Errorf("cannot marshal and write get pool response: %v", err)
		utils.HttpBadGateway(w)
	}
}
