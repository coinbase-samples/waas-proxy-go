package handlers

import (
	"context"
	"encoding/base64"

	"github.com/coinbase-samples/waas-proxy-go/config"
	log "github.com/sirupsen/logrus"

	"github.com/coinbase/waas-client-library-go/auth"
	"github.com/coinbase/waas-client-library-go/clients"
)

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
