package credentials

import (
	"net/http"

	"github.com/coinbase-samples/waas-proxy-go/config"
	"github.com/coinbase-samples/waas-proxy-go/utils"
	log "github.com/sirupsen/logrus"
)

type Credentials struct {
	ApiKeyName    string `json:"api_key_name"`
	ApiPrivateKey string `json:"api_private_key"`
}

func Retrieve(w http.ResponseWriter, r *http.Request) {

	conf := config.Get()

	creds := &Credentials{
		ApiKeyName:    conf.ApiKeyName,
		ApiPrivateKey: conf.ApiPrivateKey,
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, creds); err != nil {
		log.Errorf("Cannot marshal and write credentials response: %v", err)
		utils.HttpBadGateway(w)
	}
}
