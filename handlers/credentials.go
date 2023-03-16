package handlers

import (
	"encoding/json"
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

	body, err := json.Marshal(credentials)
	if err != nil {
		log.Errorf("Unable to marshal credentials: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatusOk(w, body); err != nil {
		log.Errorf("Cannot write pool response: %v", err)
		httpBadGateway(w)
		return
	}
}
