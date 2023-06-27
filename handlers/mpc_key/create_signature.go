package mpc_key

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/coinbase-samples/waas-proxy-go/models"
	"github.com/coinbase-samples/waas-proxy-go/utils"
	"github.com/coinbase-samples/waas-proxy-go/waas"
	v1mpckeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	common "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

const (
	ethDataMessageHashPrefix = "\x19Ethereum Signed Message:\n"
)

type SignatureRequest struct {
	Payload      string `json:"payload,omitempty"`
	PersonalSign bool   `json:"personalSign,omitempty"`
}

type SignatureResponse struct {
	Operation       string               `json:"operation,omitempty"`
	DeviceGroup     string               `json:"deviceGroup,omitempty"`
	Payload         string               `json:"payload,omitempty"`
	Signature       *v1mpckeys.Signature `json:"signature,omitempty"`
	RawTransaction  string               `json:"rawTransaction,omitempty"`
	TransactionHash string               `json:"transactionHash,omitempty"`
}

type WaitSignatureRequest struct {
	Operation   string                  `json:"operation,omitempty"`
	Transaction models.TransactionInput `json:"transaction,omitempty"`
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

	signReq := &SignatureRequest{}
	if err := json.Unmarshal(body, signReq); err != nil {
		log.Errorf("Unable to unmarshal CreateSignature request: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	payload := []byte(signReq.Payload)

	if signReq.PersonalSign {
		payloadData := string(signReq.Payload)
		generalMessage := fmt.Sprintf("%s%s", ethDataMessageHashPrefix, strconv.Itoa(len(payloadData)))
		completePayload := fmt.Sprintf("%s%s", generalMessage, payloadData)

		log.Debugf("completePayload: %s", completePayload)
		payload = []byte(completePayload)
	}

	log.Debugf("payload: %v", string(body))
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

	log.Debugf("raw create signature response: %v", response)

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

	payload := hex.EncodeToString(newSignature.GetPayload())

	var signatureRes []byte
	signatureRes = append(signatureRes, newSignature.GetSignature().GetEcdsaSignature().GetR()...)
	signatureRes = append(signatureRes, newSignature.GetSignature().GetEcdsaSignature().GetS()...)
	signatureRes = append(signatureRes, byte(newSignature.GetSignature().GetEcdsaSignature().GetV()))

	if len(signatureRes) != 65 {
		log.Errorf("unable to unmarshal signedPayload: %v", err)
		utils.HttpBadRequest(w)
		return
	}

	ethTx, err := parseTransaction(req.Transaction)
	if err != nil {
		log.Errorf("unable parse transaction to ethTx: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	signer := ethtypes.NewLondonSigner(ethTx.ChainId())

	// Create a new transaction with the given signature.
	signedTx, err := ethTx.WithSignature(signer, signatureRes)
	if err != nil {
		log.Errorf("unable to combine tx and signature: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	signedPayload, err := signedTx.MarshalBinary()
	if err != nil {
		log.Errorf("unable to marshal signed tx: %v", err)
		utils.HttpBadGateway(w)
		return
	}

	rawTransaction := hex.EncodeToString(signedPayload)

	wallet := &SignatureResponse{
		Operation:       req.Operation,
		DeviceGroup:     meta.GetDeviceGroup(),
		Payload:         payload,
		Signature:       newSignature,
		RawTransaction:  rawTransaction,
		TransactionHash: signedTx.Hash().String(),
	}

	if err := utils.HttpMarshalAndWriteJsonResponseWithOk(w, wallet); err != nil {
		log.Errorf("Cannot marshal and write create signature response: %v", err)
		utils.HttpBadGateway(w)
	}
}

// parseTransaction parses a transaction in bytes to return Transaction, Ethereum Transaction and the hashed payload
// to sign from the given serialized JSON transaction object.
func parseTransaction(tx models.TransactionInput) (*ethtypes.Transaction, error) {

	chainID, ok := new(big.Int).SetString(tx.ChainId, 0)
	if !ok {
		return nil, fmt.Errorf("invalid chainID %s", tx.ChainId)
	}

	gasTipCap, ok := new(big.Int).SetString(tx.MaxPriorityFeePerGas, 0)
	if !ok {
		return nil, fmt.Errorf("invalid maxPriorityFeePerGas %s", tx.MaxPriorityFeePerGas)
	}

	gasFeeCap, ok := new(big.Int).SetString(tx.MaxFeePerGas, 0)
	if !ok {
		return nil, fmt.Errorf("invalid maxFeePerGas %s", tx.MaxFeePerGas)
	}
	// EIP1159 uses toAddress and fromAddress, default to is 0x0000..000
	toAddress := common.HexToAddress(tx.To)

	log.Debugf("toAddress: %s - %s", toAddress, tx.ToAddress)
	value, ok := new(big.Int).SetString(tx.Value, 0)
	if !ok {
		return nil, fmt.Errorf("invalid value %s", tx.Value)
	}

	data, err := hex.DecodeString(tx.Data)
	if err != nil {
		return nil, err
	}

	accessList := ethtypes.AccessList{}

	ethTxData := &ethtypes.DynamicFeeTx{
		ChainID:    chainID,
		Nonce:      uint64(tx.Nonce),
		GasTipCap:  gasTipCap,
		GasFeeCap:  gasFeeCap,
		Gas:        uint64(tx.Gas),
		To:         &toAddress,
		Value:      value,
		Data:       data,
		AccessList: accessList,
	}

	log.Debugf("ethTxData: %v", ethTxData)
	ethTx := ethtypes.NewTx(ethTxData)
	log.Debugf("To address: %v", ethTx.To())

	return ethTx, nil
}
