package protocol

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	v1types "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/types/v1"
	log "github.com/sirupsen/logrus"
)

type broadcastRequest struct {
	NetworkId         string `json:"network,omitempty"`
	SignedTransaction string `json:"signedTransaction,omitempty"`
}

func BroadcastTransaction(w http.ResponseWriter, r *http.Request) {

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	br := &broadcastRequest{}
	if err := json.Unmarshal(body, br); err != nil {
		log.Errorf("unable to unmarshal Broadcast request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	log.Debugf("original signed: %s", br.SignedTransaction)
	encoded, err := hex.DecodeString(br.SignedTransaction)
	if err != nil {
		log.Errorf("issue with decoding hex string: %v", err)
	}

	log.Debugf("encoded tx: %s", encoded)
	req := &v1protocols.BroadcastTransactionRequest{
		Network: br.NetworkId,
		Transaction: &v1types.Transaction{
			RawSignedTransaction: encoded,
		},
	}

	tx, err := waas.GetClients().ProtocolService.BroadcastTransaction(r.Context(), req)
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
