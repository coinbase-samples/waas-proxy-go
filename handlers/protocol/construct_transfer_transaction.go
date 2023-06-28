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
	"encoding/json"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/convert"
	models "github.com/coinbase-samples/waas-proxy-go/models"
	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"
)

func ConstructTransferTransaction(w http.ResponseWriter, r *http.Request) {

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
		return
	}

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	ethInput := &models.TransactionInput{}
	if err := json.Unmarshal(body, ethInput); err != nil {
		log.Errorf("unable to unmarshal ConstructTransaction request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	finalInput, err := convert.ConvertTransferTransaction(ethInput)
	if err != nil {
		log.Errorf("cannot construct transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("sending construct transaction: %v", finalInput)
	tx, err := waas.GetClients().ProtocolService.ConstructTransferTransaction(r.Context(), finalInput)
	if err != nil {
		log.Errorf("cannot construct transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("construct transaction result: %v", tx)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("cannot marshal and wite construct tx response: %v", err)
		utils.HttpBadGateway(w)
	}

}
