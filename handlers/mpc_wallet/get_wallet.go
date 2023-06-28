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
package mpc_wallet

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	log "github.com/sirupsen/logrus"
)

func GetWallet(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	walletId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcWalletId")
	if len(walletId) == 0 {
		return
	}

	req := &v1mpcwallets.GetMPCWalletRequest{
		Name: fmt.Sprintf("pools/%s/mpcWallets/%s", poolId, walletId),
	}

	log.Debugf("getting wallet: %v", req)

	resp, err := waas.GetClients().MpcWalletService.GetMPCWallet(r.Context(), req)
	if err != nil {
		log.Errorf("cannot get wallet: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("get wallet raw response: %v", resp)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, resp); err != nil {
		log.Errorf("cannot marshal and write get wallet response: %v", err)
		utils.HttpBadGateway(w)
	}

}
