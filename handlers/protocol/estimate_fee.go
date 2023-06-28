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
package protocol

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	log "github.com/sirupsen/logrus"
)

func EstimateFee(w http.ResponseWriter, r *http.Request) {

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
		return
	}

	req := &v1protocols.EstimateFeeRequest{
		Network: fmt.Sprintf("networks/%s", networkId),
	}

	log.Debugf("sending estimageFee: %v", req)
	tx, err := waas.GetClients().ProtocolService.EstimateFee(r.Context(), req)
	if err != nil {
		log.Errorf("cannot estimateFee: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("estimateFee result: %v", tx)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("cannot marshal and wite construct tx response: %v", err)
		utils.HttpBadGateway(w)
	}

}
