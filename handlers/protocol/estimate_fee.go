package protocol

import (
	"fmt"
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	log "github.com/sirupsen/logrus"
)

func EstimateFee(w http.ResponseWriter, r *http.Request) {

	networkId := utils.HttpPathVarOrSendBadRequest(w, r, "networkId")
	if len(networkId) == 0 {
		return
	}

	req := &v1protocols.EstimateFeeRequest{
		Network: fmt.Sprintf("networks/%s", networkId),
	}

	log.Debugf("sending estimageFee: %v", req)
	tx, err := waas.GetClients().ProtocolService.EstimateFee(r.Context(), req)
	if err != nil {
		log.Errorf("cannot estimateFee: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("estimateFee result: %v", tx)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, tx); err != nil {
		log.Errorf("cannot marshal and wite construct tx response: %v", err)
		utils.HttpBadGateway(w)
	}

}
