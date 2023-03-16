package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	poolspb "github.com/WaaS-Private-Preview-v1/waas-client-library/go/coinbase/cloud/pools/v1alpha1"

	"github.com/WaaS-Private-Preview-v1/waas-client-library/go/coinbase/cloud/clients"
)

var poolServiceClient *clients.PoolServiceClient

func InitPoolClient(ctx context.Context) (err error) {

	endpoint, opts := waasClientDefaults("waas/pools")

	if poolServiceClient, err = clients.NewV1Alpha1PoolServiceClient(
		ctx,
		endpoint,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS pool client: %w", err)
	}
	return
}

func CreatePool(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read CreatePool request body: %v", err)
		httpStatusGatewayTimeout(w)
		return
	}

	req := &poolspb.CreatePoolRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal CreatePool request: %v", err)
		httpBadRequest(w)
		return
	}

	pool, err := poolServiceClient.CreatePool(ctx, req)
	if err != nil {
		log.Errorf("Cannot create pool: %v", err)
		httpBadGateway(w)
		return
	}

	body, err = json.Marshal(pool)
	if err != nil {
		log.Errorf("Cannot marshal pool struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponsWithStatusCreated(w, string(body)); err != nil {
		log.Errorf("Cannot write pool response: %v", err)
		httpBadGateway(w)
		return
	}
}
