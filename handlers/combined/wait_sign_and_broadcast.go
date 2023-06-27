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
package combined

import (
	"encoding/json"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/models"
	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	v1types "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/types/v1"
	log "github.com/sirupsen/logrus"
)

type WaitSignAndBroadcastRequest struct {
	Operation   string                  `json:"operation,omitempty"`
	Transaction models.TransactionInput `json:"transaction,omitempty"`
}

func WaitSignAndBroadcast(w http.ResponseWriter, r *http.Request) {

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	req := &WaitSignAndBroadcastRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal WaitSignAndBroadcastRequest request: %v", err)
		utils.HttpBadRequest(w)
		return
	}
	log.Debugf("waiting signature: %v", req)

	resp := waas.GetClients().MpcKeyService.CreateSignatureOperation(req.Operation)

	newSignature, err := resp.Wait(r.Context())
	if err != nil {
		log.Errorf("cannot wait create signature response: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("completed signature: %v", newSignature)

	transaction, err := convertTransaction(r.Context(), &req.Transaction)
	if err != nil {
		log.Errorf("cannot convert transaction to broadcast: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	transaction.RequiredSignatures = []*v1types.RequiredSignature{
		{
			Signature: newSignature.Signature,
		}}

	broadcastRequest := &v1protocols.BroadcastTransactionRequest{
		Network:     req.Transaction.Network,
		Transaction: transaction,
	}
	log.Debugf("broadcast request: %v", broadcastRequest)

	tx, err := waas.GetClients().ProtocolService.BroadcastTransaction(r.Context(), broadcastRequest)
	if err != nil {
		log.Errorf("cannot broadcast tx: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("broadcast result: %v", tx)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("cannot marshal and wite broadcast tx response: %v", err)
		utils.HttpBadGateway(w)
	}
}
