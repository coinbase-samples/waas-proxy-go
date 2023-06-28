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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

const (
	ethDataMessageHashPrefix = "\x19Ethereum Signed Message:\n"
	signatureLength          = 65
)

type SignatureRequest struct {
	Payload      string `json:"payload,omitempty"`
	PersonalSign bool   `json:"personalSign,omitempty"`
}

type SignatureResponse struct {
	Operation   string `json:"operation,omitempty"`
	DeviceGroup string `json:"deviceGroup,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

func CreateSignature(w http.ResponseWriter, r *http.Request) {

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

	parent := fmt.Sprintf("pools/%s/deviceGroups/%s/mpcKeys/%s", poolId, deviceGroupId, mpcKeyId)

	log.Debugf("parent: %s, body: %v", parent, string(body))

	signReq := &SignatureRequest{}
	if err := json.Unmarshal(body, signReq); err != nil {
		log.Debugf("unable to unmarshal CreateSignature request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	payload := []byte(signReq.Payload)

	if signReq.PersonalSign {
		payloadData := string(signReq.Payload)
		generalMessage := fmt.Sprintf("%s%s", ethDataMessageHashPrefix, strconv.Itoa(len(payloadData)))
		completePayload := fmt.Sprintf("%s%s", generalMessage, payloadData)

		log.Debugf("completePayload: %s", completePayload)
		payload = []byte(completePayload)
	}

	log.Debugf("payload: %v", string(payload))
	req := &v1mpckeys.CreateSignatureRequest{
		Parent: parent,
		Signature: &v1mpckeys.Signature{
			Payload: crypto.Keccak256(payload),
		},
	}

	resp, err := waas.GetClients().MpcKeyService.CreateSignature(r.Context(), req)
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
