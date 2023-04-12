package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func ConstructTransaction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	networkId, found := vars["networkId"]

	if !found {
		log.Error("Network id not passed to ConstructTransaction")
		utils.HttpBadRequest(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read ConstructTransaction request body: %v", err)
		utils.HttpGatewayTimeout(w)
		return
	}

	req := &v1protocols.ConstructTransactionRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal ConstructTransaction request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	req.Network = fmt.Sprintf("networks/%s", networkId)

	tx, err := waas.GetClients().ProtocolService.ConstructTransaction(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot construct transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("Cannot marshal and wite construct tx response: %v", err)
		utils.HttpBadGateway(w)
	}
}
