package waas

import (
	poolspb "github.com/WaaS-Private-Preview-v1/waas-client-library/go/coinbase/cloud/pools/v1alpha1"
)

type PoolService interface {
	CreatePool(ctx context.Contxt, request *poolspb.CreatePoolRequest) (*poolspb.Pool, error)
	GetPool(ctx context.Contxt, request *poolspb.GetPoolRequest) (*poolspb.Pool, error)
	Close() error
}
