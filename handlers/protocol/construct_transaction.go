package protocol

import (
	"encoding/json"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/convert"
	models "github.com/coinbase-samples/waas-proxy-go/models"
	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	log "github.com/sirupsen/logrus"
)

func ConstructTransaction(w http.ResponseWriter, r *http.Request) {

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
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

	req, err := convert.ConvertEip1559Transaction(ethInput)
	if err != nil {
		log.Errorf("cannot construct transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("sending construct transaction: %v", req)
	tx, err := waas.GetClients().ProtocolService.ConstructTransaction(r.Context(), req)
	if err != nil {
		log.Errorf("cannot construct transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("construct transaction result: %v", tx)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("cannot marshal and wite construct tx response: %v", err)
		utils.HttpBadGateway(w)
	}

}
