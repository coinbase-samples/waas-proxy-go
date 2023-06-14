package mpc_transaction

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpctransactions "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_transactions/v1"
	log "github.com/sirupsen/logrus"
)

func GetMpcTransaction(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	walletId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcWalletId")
	if len(walletId) == 0 {
		return
	}

	transactionId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcTransactionId")
	if len(transactionId) == 0 {
		return
	}

	req := &v1mpctransactions.GetMPCTransactionRequest{
		Name: fmt.Sprintf("pools/%s/mpcWallets/%s/mpcTransactions/%s", poolId, walletId, transactionId),
	}

	log.Debugf("getting mpc transactions: %v", req)
	response, err := waas.GetClients().MpcTransactionService.GetMPCTransaction(r.Context(), req)

	if err != nil {
		log.Errorf("Cannot get transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("get transaction raw response: %v", response)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write get mpc wallet response: %v", err)
		utils.HttpBadGateway(w)
	}
}
