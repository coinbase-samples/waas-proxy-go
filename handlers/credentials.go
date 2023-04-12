package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Credentials struct {
	ApiKeyName    string `json:"api_key_name"`
	ApiPrivateKey string `json:"api_private_key"`
}

func RetrieveCredentials(w http.ResponseWriter, r *http.Request) {

	credentials := &Credentials{
		ApiKeyName:    appConfig.ApiKeyName,
		ApiPrivateKey: appConfig.ApiPrivateKey,
	}

	if err := marhsallAndWriteJsonResponseWithOk(w, credentials); err != nil {
		log.Errorf("Cannot marshal and write credentials response: %v", err)
		httpBadGateway(w)
	}
}
