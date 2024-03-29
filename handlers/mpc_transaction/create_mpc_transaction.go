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
package mpc_transaction

import (
	"fmt"
	"net/http"
	"time"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpctransactions "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_transactions/v1"
	v1 "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/types/v1"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

func CreateMPCTransaction(w http.ResponseWriter, r *http.Request) {
	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	mpcWalletId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcWalletId")
	if len(mpcWalletId) == 0 {
		return
	}

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	log.Debugf("raw body: %s", string(body))
	input := &v1.Transaction{}
	if err := protojson.Unmarshal(body, input); err != nil {
		log.Errorf("unable to unmarshal CreateMPCTransaction request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	fromAddress := input.Input.GetEthereum_1559Input()
	log.Debugf("1559 input: %v", input)
	log.Debugf("fromAddress: %s", fromAddress.FromAddress)
	req := &v1mpctransactions.CreateMPCTransactionRequest{
		Parent:    fmt.Sprintf("pools/%s/mpcWallets/%s", poolId, mpcWalletId),
		RequestId: uuid.New().String(),
		Input:     input.Input,
		MpcTransaction: &v1mpctransactions.MPCTransaction{
			Network:       "networks/ethereum-goerli",
			FromAddresses: []string{fmt.Sprintf("networks/ethereum-goerli/addresses/%s", fromAddress.FromAddress)},
		},
	}

	log.Debugf("creating mpc transaction: %v", req)

	resp, err := waas.GetClients().MpcTransactionService.CreateMPCTransaction(r.Context(), req)
	if err != nil {
		log.Errorf("cannot create mpc transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	var metadata *v1mpctransactions.CreateMPCTransactionMetadata

	for metadata == nil || metadata.GetDeviceGroup() == "" {
		time.Sleep(20 * time.Millisecond)
		metadata, err = resp.Metadata()
		log.Debugf("new metadata: %v", metadata)
	}

	log.Debugf("create mpc result: %v", metadata)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, metadata); err != nil {
		log.Errorf("cannot marshal and write create mpc transaction response: %v", err)
		utils.HttpBadGateway(w)
	}

}
