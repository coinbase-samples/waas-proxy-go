package mpc_wallet

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpcwallets "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

func ListWallets(w http.ResponseWriter, r *http.Request) {

	// TODO: This needs to page for the end client - iterator blasts through everything

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	params, _ := url.ParseQuery(r.URL.RawQuery)
	filterDeviceGroup := params.Get("deviceGroup")
	log.Debugf("filtering by deviceGroup: %s", filterDeviceGroup)
	wallets := listWallet(r.Context(), poolId, filterDeviceGroup)

	log.Debugf("found wallets %d", len(wallets))
	start := time.Now().Unix()
	if len(wallets) == 0 {
		for i := 0; i < 60; i++ {
			time.Sleep(time.Second)
			log.Debugf("slept, fetching again: %d", time.Now().Unix())
			wallets = listWallet(r.Context(), poolId, filterDeviceGroup)
			if len(wallets) > 0 {
				break
			}
		}
	}
	end := time.Now().Unix()

	log.Debugf("time to find wallet: %d", end-start)

	response := &v1mpcwallets.ListMPCWalletsResponse{MpcWallets: wallets}

	log.Debugf("returning wallets: %v", response)
	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write mpc wallet list response: %v", err)
		utils.HttpBadGateway(w)
	}
}

func listWallet(ctx context.Context, poolId, deviceGroup string) []*v1mpcwallets.MPCWallet {
	var wallets []*v1mpcwallets.MPCWallet

	parent := fmt.Sprintf("pools/%s", poolId)
	log.Debugf("listing wallets under parent: %s", parent)
	req := &v1mpcwallets.ListMPCWalletsRequest{
		Parent:   parent,
		PageSize: 100,
	}
	count := 0
	iter := waas.GetClients().MpcWalletService.ListMPCWallets(ctx, req)
	for {
		wallet, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Errorf("Cannot iterate wallets: %v", err)
			return wallets
		}
		if wallet.DeviceGroup == deviceGroup || deviceGroup == "" {
			wallets = append(wallets, wallet)
		}
		count++
	}
	log.Debugf("checked %d wallets, found %d", count, len(wallets))
	return wallets
}
