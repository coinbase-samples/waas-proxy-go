package handlers

type CreatePoolRequest struct {
	DisplayName string `json:"displayName"`
	PoolId      string `json:"poolId"`
}
