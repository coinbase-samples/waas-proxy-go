package handlers

import (
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/WaaS-Private-Preview-v1/waas-client-library/go/coinbase/cloud/clients"
)

const defaultWaaSApiHost = "https://cloud-api-beta.coinbase.com"

func waasClientDefaults(endpointPath string) (endpoint string, opts []clients.Option) {

	url, err := url.Parse(defaultWaaSApiHost)
	if err != nil {
		log.Fatalf("Invalid host: %s - %v", defaultWaaSApiHost, err)
	}

	url.Path = endpointPath

	endpoint = url.String()

	// Looks for env vars: COINBASE_CLOUD_API_KEY_NAME, COINBASE_CLOUD_API_PRIVATE_KEY
	opts = []clients.Option{
		clients.WithCloudAPIKeyAuth(clients.WithDefaultENVVariables()),
	}

	return

}
