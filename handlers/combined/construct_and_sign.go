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
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/convert"
	"github.com/coinbase-samples/waas-proxy-go/models"
	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	v1types "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/types/v1"
	log "github.com/sirupsen/logrus"
)

type SignatureResponse struct {
	Operation   string `json:"operation,omitempty"`
	DeviceGroup string `json:"deviceGroup,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

func ConstructAndSign(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	deviceGroupId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceGroupId")
	if len(deviceGroupId) == 0 {
		return
	}

	mpcKeyId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcKeyId")
	if len(mpcKeyId) == 0 {
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

	tx, err := convertTransaction(r.Context(), ethInput)
	if err != nil {
		log.Errorf("cannot convert transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	parent := fmt.Sprintf("pools/%s/deviceGroups/%s/mpcKeys/%s", poolId, deviceGroupId, mpcKeyId)

	log.Debugf("parent: %s, body: %v", parent, string(body))

	log.Debugf("payload: %v", string(body))
	createSigReq := &v1mpckeys.CreateSignatureRequest{
		Parent: parent,
		Signature: &v1mpckeys.Signature{
			Payload: tx.RequiredSignatures[0].Payload,
		},
	}

	resp, err := waas.GetClients().MpcKeyService.CreateSignature(r.Context(), createSigReq)
	if err != nil {
		log.Errorf("cannot createSignature operations: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("createSig response: %v", resp)
	meta, err := resp.Metadata()
	if err != nil {
		log.Errorf("cannot metadata createSignature: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	response := &SignatureResponse{
		Operation:   resp.Name(),
		DeviceGroup: meta.GetDeviceGroup(),
		Payload:     string(meta.GetPayload()),
	}

	log.Debugf("raw create signature response: %v", response)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("cannot marshal and write create signature metadata response: %v", err)
		utils.HttpBadGateway(w)
	}
}

func convertTransaction(ctx context.Context, ethInput *models.TransactionInput) (*v1types.Transaction, error) {
	var tx *v1types.Transaction

	if len(ethInput.Asset) > 0 {
		req, err := convert.ConvertTransferTransaction(ethInput)
		if err != nil {
			return tx, err
		}

		log.Debugf("sending construct transfer transaction: %v", req)
		tx, err = waas.GetClients().ProtocolService.ConstructTransferTransaction(ctx, req)
		if err != nil {
			return tx, err
		}
	} else {
		req, err := convert.ConvertEip1559Transaction(ethInput)
		if err != nil {
			return tx, err
		}

		log.Debugf("sending construct transaction: %v", req)
		tx, err = waas.GetClients().ProtocolService.ConstructTransaction(ctx, req)
		if err != nil {
			return tx, err
		}
	}
	return tx, nil
}
