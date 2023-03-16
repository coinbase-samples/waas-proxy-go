package handlers

import (
	"context"
	"encoding/base64"
	"net/url"

	"github.com/coinbase-samples/waas-proxy-go/config"
	log "github.com/sirupsen/logrus"

	"github.com/WaaS-Private-Preview-v1/waas-client-library/go/coinbase/cloud/clients"
)

const defaultWaaSApiHost = "https://cloud-api-beta.coinbase.com"

func initWaaSClients(config config.AppConfig) error {

	if err := initPoolClient(context.Background(), config); err != nil {
		return err
	}

	return nil
}

func waasClientDefaults(
	config config.AppConfig,
	endpointPath string,
) (endpoint string, opts []clients.Option) {

	url, err := url.Parse(defaultWaaSApiHost)
	if err != nil {
		log.Fatalf("Invalid host: %s - %v", defaultWaaSApiHost, err)
	}

	url.Path = endpointPath

	endpoint = url.String()

	apiPrivateKey, err := base64.StdEncoding.DecodeString(config.ApiPrivateKey)
	if err != nil {
		log.Fatalf("Cannot base64 decode private key: %v", err)
	}

	opts = []clients.Option{
		clients.WithCloudAPIKeyAuth(clients.WithAPIKey(config.ApiKeyName, string(apiPrivateKey))),
	}

	return
}
