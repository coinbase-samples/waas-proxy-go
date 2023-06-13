package protocol

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/convert"
	models "github.com/coinbase-samples/waas-proxy-go/models"
	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	v1 "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/types/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
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
	if err := protojson.Unmarshal(body, ethInput); err != nil {
		log.Errorf("Unable to unmarshal ConstructTransaction request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	finalInput, err := convert.ConvertEip1559Transaction(ethInput)
	if err != nil {
		log.Errorf("Cannot construct transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	input := &v1.TransactionInput{
		Input: &v1.TransactionInput_Ethereum_1559Input{
			Ethereum_1559Input: finalInput,
		},
	}
	req := &v1protocols.ConstructTransactionRequest{
		Input:   input,
		Network: fmt.Sprintf("networks/%s", networkId),
	}

	log.Debugf("sending construct transaction: %v", req)
	tx, err := waas.GetClients().ProtocolService.ConstructTransaction(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot construct transaction: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("construct transaction result: %v", tx)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("Cannot marshal and wite construct tx response: %v", err)
		utils.HttpBadGateway(w)
	}

}
