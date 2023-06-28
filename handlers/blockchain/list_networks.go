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
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"

	v1blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
)

func ListNetworks(w http.ResponseWriter, r *http.Request) {

	req := &v1blockchain.ListNetworksRequest{}

	log.Debugf("ListNetworks request: %v", req)
	requestPageInfo, err := utils.HttpRequestPageInfo(r)
	if err != nil {
		utils.HttpBadRequest(w)
		return
	}

	if requestPageInfo.Passed() {
		req.PageToken = requestPageInfo.Token
		req.PageSize = requestPageInfo.Size
	}

	iter := waas.GetClients().BlockchainService.ListNetworks(r.Context(), req)

	_, err = iter.Next()
	if err != nil && err != iterator.Done {
		log.Errorf("cannot iterate assets: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	response := iter.Response()

	log.Debugf("ListNetworks response: %v", response)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("cannot marshal and write list networks response: %v", err)
		utils.HttpBadGateway(w)
	}
}
