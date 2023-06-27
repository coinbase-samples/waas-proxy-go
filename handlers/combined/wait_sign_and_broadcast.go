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
		log.Errorf("Cannot wait create signature response: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("completed signature: %v", newSignature)

	transaction, err := convertTransaction(r.Context(), &req.Transaction)
	if err != nil {
		log.Errorf("Cannot convert transaction to broadcast: %v", err)
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
		log.Errorf("Cannot broadcast tx: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("broadcast result: %v", tx)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("Cannot marshal and wite broadcast tx response: %v", err)
		utils.HttpBadGateway(w)
	}
}
