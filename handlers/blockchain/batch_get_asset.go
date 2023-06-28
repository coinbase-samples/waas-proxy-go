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
package blockchain

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	log "github.com/sirupsen/logrus"
)

func BatchGetAsset(w http.ResponseWriter, r *http.Request, networkId string, names []string) {

	req := &v1blockchain.BatchGetAssetsRequest{
		Parent: fmt.Sprintf("networks/%s", networkId),
		Names:  names,
	}

	log.Debugf("BatchGetAssets request: %v", req)
	resp, err := waas.GetClients().BlockchainService.BatchGetAssets(r.Context(), req)
	if err != nil {
		utils.HttpBadGateway(w)
		log.Error(err)
		return
	}

	log.Debugf("BatchGetAssets response: %v", resp)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("cannot marshal and write get device group response: %v", err)
		utils.HttpBadGateway(w)
	}
}
