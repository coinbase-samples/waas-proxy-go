package models

type TransactionInput struct {
	ChainId string `json:"chainId,omitempty"`
	// The nonce of the transaction. This value may be ignored depending on the API.
	Nonce uint64 `json:"nonce,omitempty"`
	// The EIP-1559 maximum priority fee per gas either as a "0x"-prefixed hex string or a base-10 number.
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas,omitempty"`
	// The EIP-1559 maximum fee per gas either as a "0x"-prefixed hex string or a base-10 number.
	MaxFeePerGas string `json:"maxFeePerGas,omitempty"`
	// The maximum amount of gas to use on the transaction.
	Gas uint64 `json:"gas,omitempty"`
	// The checksummed address from which the transaction will originate, as a "0x"-prefixed hex string.
	// Note: This is NOT a WaaS Address resource of the form
	// networks/{networkID}/addresses/{addressID}.
	FromAddress string `json:"fromAddress,omitempty"`
	// The checksummed address to which the transaction is addressed, as a "0x"-prefixed hex string.
	// Note: This is NOT a WaaS Address resource of the form
	// networks/{networkID}/addresses/{addressID}.
	ToAddress string `json:"toAddress,omitempty"`
	// The native value of the transaction as a "0x"-prefixed hex string or a base-10 number.
	Value string `json:"value,omitempty"`
	// The data for the transaction.
	Data string `json:"data,omitempty"`
}
