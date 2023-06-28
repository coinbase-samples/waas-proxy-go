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

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpctransactions "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_transactions/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

func ListMpcTransactions(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	mpcWalletId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcWalletId")
	if len(mpcWalletId) == 0 {
		return
	}

	req := &v1mpctransactions.ListMPCTransactionsRequest{
		Parent: fmt.Sprintf("pools/%s/mpcWallets/%s", poolId, mpcWalletId),
	}

	log.Debugf("listing mpc transactions: %v", req)
	iter := waas.GetClients().MpcTransactionService.ListMPCTransactions(r.Context(), req)

	var transactions []*v1mpctransactions.MPCTransaction
	for {
		transaction, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("cannot iterate transactions: %v", err)
			utils.HttpBadGateway(w)
			return
		}
		transactions = append(transactions, transaction)
	}

	response := &v1mpctransactions.ListMPCTransactionsResponse{MpcTransactions: transactions}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("cannot marshal and write mpc wallet list response: %v", err)
		utils.HttpBadGateway(w)
	}
}
