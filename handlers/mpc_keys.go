package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/coinbase-samples/waas-proxy-go/config"
	waasv1 "github.com/coinbase/waas-client-library-go/clients/v1"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	ethDataMessageHashPrefix = "\x19Ethereum Signed Message:\n"
)

type CreateDeviceGroupResponse struct {
	// The resource name of the Balance.
	// Format: operations/{operation_id}
	Operation   string `json:"operation,omitempty"`
	DeviceGroup string `json:"deviceGroup,omitempty"`
}

var mpcKeysServiceClient *waasv1.MPCKeyServiceClient

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

func MpcWalletListOperations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Pool id not passed to MpcWalletListOperations")
		httpBadRequest(w)
		return
	}

	deviceGroupId, found := vars["deviceGroupId"]
	if !found {
		log.Error("Device Group Id not passed to MpcWalletListOperations")
		httpBadRequest(w)
		return
	}

	req := &v1mpckeys.ListMPCOperationsRequest{
		Parent: fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId),
	}

	log.Infof("listing mpc op request: %v", req)
	resp, err := mpcKeysServiceClient.ListMPCOperations(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot list mpc operations: %v", err)
		httpBadGateway(w)
		return
	}
	body, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Cannot marshal tx struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatus(w, body, http.StatusOK); err != nil {
		log.Errorf("Cannot write tx response: %v", err)
		httpBadGateway(w)
		return
	}

}

func MpcCreateDeviceGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Pool id not passed to MpcWalletCreate")
		httpBadRequest(w)
		return
	}

	deviceId, found := vars["deviceId"]
	if !found {
		log.Error("Device Id not passed to MpcWalletListOperations")
		httpBadRequest(w)
		return
	}

	req := &v1mpckeys.CreateDeviceGroupRequest{
		Parent: fmt.Sprintf("pools/%s/device/%s", poolId, deviceId),
	}

	resp, err := mpcKeysServiceClient.CreateDeviceGroup(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot list mpc operations: %v", err)
		httpBadGateway(w)
		return
	}
	metadata, _ := resp.Metadata()

	finalResp := &CreateDeviceGroupResponse{
		Operation:   resp.Name(),
		DeviceGroup: metadata.GetDeviceGroup(),
	}
	body, err := json.Marshal(finalResp)
	if err != nil {
		log.Errorf("Cannot marshal tx struct: %v", err)
		httpBadGateway(w)
		return
	}

	if err = writeJsonResponseWithStatus(w, body, http.StatusOK); err != nil {
		log.Errorf("Cannot write tx response: %v", err)
		httpBadGateway(w)
		return
	}

}

func MpcRegisterDevice(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read RegisterDevice request body: %v", err)
		httpGatewayTimeout(w)
		return
	}

	req := &v1mpckeys.RegisterDeviceRequest{}
	if err := json.Unmarshal(body, req); err != nil {
		log.Errorf("Unable to unmarshal RegisterDevice request: %v", err)
		httpBadRequest(w)
		return
	}
	log.Infof("registering device: %v", req)

	resp, err := mpcKeysServiceClient.RegisterDevice(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot register new device: %v", err)
		httpBadGateway(w)
		return
	}
	log.Infof("register device raw response: %v", resp)
	body, err = json.Marshal(resp)
	if err != nil {
		log.Errorf("Cannot marshal register device struct: %v", err)
		httpBadGateway(w)
		return
	}

	log.Infof("register device result: %v", string(body))

	if err = writeJsonResponseWithStatus(w, body, http.StatusOK); err != nil {
		log.Errorf("Cannot write register device response: %v", err)
		httpBadGateway(w)
		return
	}
}

func MpcCreateSignature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	poolId, found := vars["poolId"]
	if !found {
		log.Error("Pool id not passed to MpcWalletCreate")
		httpBadRequest(w)
		return
	}

	deviceGroupId, found := vars["deviceGroupId"]
	if !found {
		log.Error("Device Group Id not passed to MpcCreateSignature")
		httpBadRequest(w)
		return
	}
	mpcKeyId, found := vars["mpcKeyId"]
	if !found {
		log.Error("MpcKeyId Id not passed to MpcCreateSignature")
		httpBadRequest(w)
		return
	}

	parent := fmt.Sprintf("pools/%s/deviceGroups/%s/mpcKeys/%s", poolId, deviceGroupId, mpcKeyId)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Unable to read RegisterDevice request body: %v", err)
		httpGatewayTimeout(w)
		return
	}

	log.Infof("parent: %s, body: %v", parent, string(body))

	payloadData := string(body)
	generalMessage := fmt.Sprintf("%s%s", "\x19Ethereum Signed Message:\n", strconv.Itoa(len(payloadData)))
	completePayload := fmt.Sprintf("%s%s", generalMessage, payloadData)

	log.Infof("completePayload: %s", completePayload)
	payload := []byte(completePayload)

	req := &v1mpckeys.CreateSignatureRequest{
		Parent: parent,
		Signature: &v1mpckeys.Signature{
			Payload: crypto.Keccak256(payload),
		},
	}
	/*
		if err := json.Unmarshal(body, req); err != nil {
			log.Errorf("Unable to unmarshal createSignature request: %v", err)
			httpBadRequest(w)
			return
		}
	*/

	resp, err := mpcKeysServiceClient.CreateSignature(r.Context(), req)
	if err != nil {
		log.Errorf("Cannot createSignature operations: %v", err)
		httpBadGateway(w)
		return
	}

	log.Infof("createSig response: %v", resp)
	sig, err := resp.Poll(r.Context())
	if err != nil {
		log.Errorf("Cannot poll createSignature: %v", err)
		httpBadGateway(w)
		return
	}

	log.Infof("after poll: %v", sig)

	mpcParent := fmt.Sprintf("pools/%s/deviceGroups/%s", poolId, deviceGroupId)
	var mpcResp *v1mpckeys.ListMPCOperationsResponse
	counter := 1
	for counter < 20 {
		log.Infof("listing mpc operations %s: %s", fmt.Sprint(counter), mpcParent)
		mpcResp, err = mpcKeysServiceClient.ListMPCOperations(r.Context(), &v1mpckeys.ListMPCOperationsRequest{
			Parent: mpcParent,
		})

		log.Infof("list mpc ops response: %v", mpcResp)
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
		httpBadGateway(w)
		return
	}

	log.Infof("raw response: %v", mpcResp)

	body, err = json.Marshal(mpcResp)
	if err != nil {
		log.Errorf("Cannot marshal createSignature metadata struct: %v", err)
		httpBadGateway(w)
		return
	}
	log.Infof("final response: %s", string(body))

	if err = writeJsonResponseWithStatus(w, body, http.StatusOK); err != nil {
		log.Errorf("Cannot write createSignature metadata response: %v", err)
		httpBadGateway(w)
		return
	}

}
