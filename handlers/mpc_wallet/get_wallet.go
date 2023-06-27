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
