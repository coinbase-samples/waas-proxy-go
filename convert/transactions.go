/**
 * Copyright 2023 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package convert

import (
	"encoding/hex"

	models "github.com/coinbase-samples/waas-proxy-go/models"
	ethereumpb "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/ethereum/v1"
	v1protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	v1 "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/types/v1"
)

func ConvertEip1559Transaction(t *models.TransactionInput) (*v1protocols.ConstructTransactionRequest, error) {
	data, err := hex.DecodeString(t.Data)
	if err != nil {
		return nil, err
	}
	ethInput := &ethereumpb.EIP1559TransactionInput{
		ChainId:              t.ChainId,
		Nonce:                t.Nonce,
		MaxPriorityFeePerGas: t.MaxPriorityFeePerGas,
		MaxFeePerGas:         t.MaxFeePerGas,
		Gas:                  t.Gas,
		FromAddress:          t.FromAddress,
		ToAddress:            t.ToAddress,
		Value:                t.Value,
		Data:                 data,
	}
	input := &v1.TransactionInput{
		Input: &v1.TransactionInput_Ethereum_1559Input{
			Ethereum_1559Input: ethInput,
		},
	}
	req := &v1protocols.ConstructTransactionRequest{
		Input:   input,
		Network: t.Network,
	}
	return req, nil
}

func ConvertTransferTransaction(t *models.TransactionInput) (*v1protocols.ConstructTransferTransactionRequest, error) {
	return &v1protocols.ConstructTransferTransactionRequest{
		Network: t.Network,
		Asset:   t.Asset,
		Nonce:   int64(t.Nonce),
		Fee: &v1.TransactionFee{
			Fee: &v1.TransactionFee_EthereumFee{
				EthereumFee: &ethereumpb.DynamicFeeInput{
					MaxPriorityFeePerGas: t.MaxPriorityFeePerGas,
					MaxFeePerGas:         t.MaxFeePerGas,
				},
			},
		},
		Sender:    t.FromAddress,
		Recipient: t.ToAddress,
		Amount:    t.Value,
	}, nil
}
