package convert

import (
	"encoding/hex"

	models "github.com/coinbase-samples/waas-proxy-go/models"
	ethereumpb "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/ethereum/v1"
)

func ConvertEip1559Transaction(t *models.TransactionInput) (*ethereumpb.EIP1559TransactionInput, error) {
	data, err := hex.DecodeString(t.Data)
	if err != nil {
		return nil, err
	}
	return &ethereumpb.EIP1559TransactionInput{
		ChainId:              t.ChainId,
		Nonce:                t.Nonce,
		MaxPriorityFeePerGas: t.MaxPriorityFeePerGas,
		MaxFeePerGas:         t.MaxFeePerGas,
		Gas:                  t.Gas,
		FromAddress:          t.FromAddress,
		ToAddress:            t.ToAddress,
		Value:                t.Value,
		Data:                 data,
	}, nil
}
