package handlers

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/coinbase-samples/waas-proxy-go/config"
	log "github.com/sirupsen/logrus"

	"github.com/coinbase/waas-client-library-go/auth"
	"github.com/coinbase/waas-client-library-go/clients"
	waasv1 "github.com/coinbase/waas-client-library-go/clients/v1"
)

var blockchainServiceClient *waasv1.BlockchainServiceClient
var mpcKeysServiceClient *waasv1.MPCKeyServiceClient
var mpcWalletServiceClient *waasv1.MPCWalletServiceClient

var poolServiceClient *waasv1.PoolServiceClient

func initWaaSClients(config config.AppConfig) error {

	if err := initBlockchainClient(context.Background(), config); err != nil {
		return err
	}

	if err := initPoolClient(context.Background(), config); err != nil {
		return err
	}

	if err := initProtocolClient(context.Background(), config); err != nil {
		return err
	}

	if err := initMpcWalletClient(context.Background(), config); err != nil {
		return err
	}

	if err := initMpcKeyClient(context.Background(), config); err != nil {
		return err
	}

	return nil
}

func initMpcKeyClient(ctx context.Context, config config.AppConfig) (err error) {

	opts := waasClientDefaults(config)

	if mpcKeysServiceClient, err = waasv1.NewMPCKeyServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS mpc key client: %w", err)
	}

	return
}

func initBlockchainClient(ctx context.Context, config config.AppConfig) (err error) {

	opts := waasClientDefaults(config)

	if blockchainServiceClient, err = waasv1.NewBlockchainServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS blockchain client: %w", err)
	}
	return
}

func initMpcWalletClient(ctx context.Context, config config.AppConfig) error {
	var e, err error
	opts := waasClientDefaults(config)

	if mpcWalletServiceClient, err = waasv1.NewMPCWalletServiceClient(
		ctx,
		opts...,
	); err != nil {
		e = fmt.Errorf("unable to init WaaS mpc wallet client: %w", err)
	}

	return e
}

func initPoolClient(ctx context.Context, config config.AppConfig) (err error) {

	opts := waasClientDefaults(config)

	if poolServiceClient, err = waasv1.NewPoolServiceClient(
		ctx,
		opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS pool client: %w", err)
	}
	return
}

func waasClientDefaults(
	config config.AppConfig,
) (opts []clients.WaaSClientOption) {

	apiPrivateKey, err := base64.StdEncoding.DecodeString(config.ApiPrivateKey)
	if err != nil {
		log.Fatalf("Cannot base64 decode private key: %v", err)
	}

	opts = []clients.WaaSClientOption{
		clients.WithAPIKey(
			&auth.APIKey{
				Name:       config.ApiKeyName,
				PrivateKey: string(apiPrivateKey),
			},
		),
	}
	return
}
