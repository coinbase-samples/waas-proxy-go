package mpc_key

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

const (
	ethDataMessageHashPrefix = "\x19Ethereum Signed Message:\n"
)

type SignatureResponse struct {
	// The resource name of the Balance.
	// Format: operations/{operation_id}
	Operation      string               `json:"operation,omitempty"`
	DeviceGroup    string               `json:"deviceGroup,omitempty"`
	Payload        string               `json:"payload,omitempty"`
	Signature      *v1mpckeys.Signature `json:"signature,omitempty"`
	RawTransaction string               `json:"rawTransaction,omitempty"`
}

type WaitSignatureRequest struct {
	Operation string `json:"operation,omitempty"`
}

func CreateSignature(w http.ResponseWriter, r *http.Request) {

	poolId := utils.HttpPathVarOrSendBadRequest(w, r, "poolId")
	if len(poolId) == 0 {
		return
	}

	deviceGroupId := utils.HttpPathVarOrSendBadRequest(w, r, "deviceGroupId")
	if len(deviceGroupId) == 0 {
		return
	}

	mpcKeyId := utils.HttpPathVarOrSendBadRequest(w, r, "mpcKeyId")
	if len(mpcKeyId) == 0 {
		return
	}

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	parent := fmt.Sprintf("pools/%s/deviceGroups/%s/mpcKeys/%s", poolId, deviceGroupId, mpcKeyId)

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
	meta, err := resp.Metadata()
	if err != nil {
		log.Errorf("Cannot metadata createSignature: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	response := &SignatureResponse{
		Operation:   resp.Name(),
		DeviceGroup: meta.GetDeviceGroup(),
		Payload:     string(meta.GetPayload()),
	}

	log.Debugf("raw response: %v", response)

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, response); err != nil {
		log.Errorf("Cannot marshal and write create signature metadata response: %v", err)
		utils.HttpBadGateway(w)
	}
}

func WaitSignature(w http.ResponseWriter, r *http.Request) {

	body, err := utils.HttpReadBodyOrSendGatewayTimeout(w, r)
	if err != nil {
		return
	}

	req := &WaitSignatureRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("unable to unmarshal WaitSignature request: %v", err)
		utils.HttpBadRequest(w)
		return
	}
	log.Debugf("waiting signature: %v", req)

	resp := waas.GetClients().MpcKeyService.CreateSignatureOperation(req.Operation)

	newSignature, err := resp.Wait(r.Context())
	if err != nil {
		log.Errorf("Cannot wait create signature response: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	log.Debugf("completed signature: %v", newSignature)
	meta, err := resp.Metadata()

	if err != nil {
		log.Errorf("Cannot get metadata create signature response: %v", err)
		utils.HttpBadGateway(w)
		return
	}
	/*
		ecdsaSig := newSignature.Signature.GetEcdsaSignature()
		rVal := ecdsaSig.GetR()
		sVal := ecdsaSig.GetS()
		vVal := ecdsaSig.GetV()
		rawTransaction := []byte{}
		rawTransaction = append(rawTransaction, rVal...)
		rawTransaction = append(rawTransaction, sVal...)
		rawTransaction = append(rawTransaction, []byte(vVal)...)

		log.Debugf("rawTransaction: %v", rawTransaction)
	*/
	wallet := &SignatureResponse{
		Operation:   req.Operation,
		DeviceGroup: meta.GetDeviceGroup(),
		Payload:     string(meta.GetPayload()),
		Signature:   newSignature,
		//RawTransaction: string(rawTransaction),
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, wallet); err != nil {
		log.Errorf("Cannot marshal and write create signature response: %v", err)
		utils.HttpBadGateway(w)
	}
}
