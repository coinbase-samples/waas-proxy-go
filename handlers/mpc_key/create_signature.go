package mpc_key

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	ethDataMessageHashPrefix = "\x19Ethereum Signed Message:\n"
)

func CreateSignature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Pool id not passed to MpcWalletCreate")
		utils.HttpBadRequest(w)
		return
	}

	deviceGroupId, found := vars["deviceGroupId"]
	if !found {
		log.Error("Device Group Id not passed to CreateSignature")
		utils.HttpBadRequest(w)
		return
	}
	mpcKeyId, found := vars["mpcKeyId"]
	if !found {
		log.Error("MpcKeyId Id not passed to CreateSignature")
		utils.HttpBadRequest(w)
		return
	}

	parent := fmt.Sprintf("pools/%s/deviceGroups/%s/mpcKeys/%s", poolId, deviceGroupId, mpcKeyId)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read RegisterDevice request body: %v", err)
		utils.HttpGatewayTimeout(w)
		return
	}

	log.Debugf("parent: %s, body: %v", parent, string(body))

	payloadData := string(body)
	generalMessage := fmt.Sprintf("%s%s", ethDataMessageHashPrefix, strconv.Itoa(len(payloadData)))
	completePayload := fmt.Sprintf("%s%s", generalMessage, payloadData)

	log.Debugf("completePayload: %s", completePayload)
	payload := []byte(completePayload)

	req := &v1mpckeys.CreateSignatureRequest{
		Parent: parent,
		Signature: &v1mpckeys.Signature{
			Payload: crypto.Keccak256(payload),
		},
	}

	resp, err := waas.GetClients().MpcKeyService.CreateSignature(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot createSignature operations: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("createSig response: %v", resp)
	sig, err := resp.Poll(r.Context())
	if err != nil {
		log.Errorf("Cannot poll createSignature: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("after poll: %v", sig)

	mpcParent := fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId)
	var mpcResp *v1mpckeys.ListMPCOperationsResponse
	counter := 1
	for counter < 20 {
		log.Debugf("listing mpc operations %s: %s", fmt.Sprint(counter), mpcParent)
		mpcResp, err = waas.GetClients().MpcKeyService.ListMPCOperations(r.Context(), &v1mpckeys.ListMPCOperationsRequest{
			Parent: mpcParent,
		})

		log.Debugf("list mpc ops response: %v", mpcResp)
		if err != nil {
			log.Errorf("cannot list mpc operations for %s: %v", mpcParent, err)
			time.Sleep(250 * time.Millisecond)
			counter = counter + 1
		}

		if mpcResp != nil && len(mpcResp.MpcOperations) > 0 {
			break
		} else {
			time.Sleep(250 * time.Millisecond)
			counter = counter + 1
		}
	}

	if err != nil {
		log.Errorf("Cannot list mpc operations: %v, parent: %s", err, mpcParent)
		utils.HttpBadGateway(w)
		return
	}

	log.Debugf("raw response: %v", mpcResp)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, mpcResp); err != nil {
		log.Errorf("Cannot marshal and write create signature metadata response: %v", err)
		utils.HttpBadGateway(w)
	}
}