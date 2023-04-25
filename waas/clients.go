package waas

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/coinbase-samples/waas-proxy-go/config"
	log "github.com/sirupsen/logrus"

	"github.com/coinbase/waas-client-library-go/auth"

	waasClients "github.com/coinbase/waas-client-library-go/clients"
	waasv1 "github.com/coinbase/waas-client-library-go/clients/v1"
)

var clients *Clients

type Clients struct {
	BlockchainService     *waasv1.BlockchainServiceClient
	MpcKeyService         *waasv1.MPCKeyServiceClient
	MpcWalletService      *waasv1.MPCWalletServiceClient
	MpcTransactionService *waasv1.MPCTransactionServiceClient
	PoolService           *waasv1.PoolServiceClient
	ProtocolService       *waasv1.ProtocolServiceClient
}

type clientInit struct {
	config config.AppConfig
	opts   []waasClients.WaaSClientOption
	c      *Clients
}

func GetClients() *Clients {
	return clients
}

func InitClients(config config.AppConfig) error {

	ctx := context.Background()

	clients = &Clients{}

	init := &clientInit{
		config: config,
		opts:   waasClientDefaults(config),
		c:      clients,
	}

	if err := initBlockchainClient(ctx, init); err != nil {
		return err
	}

	if err := initPoolClient(ctx, init); err != nil {
		return err
	}

	if err := initProtocolClient(ctx, init); err != nil {
		return err
	}

	if err := initMpcWalletClient(ctx, init); err != nil {
		return err
	}

	if err := initMpcKeyClient(ctx, init); err != nil {
		return err
	}

	if err := initMpcTransactionClient(ctx, init); err != nil {
		return err
	}

	return nil
}

// initProtocolClient creates a new protocol service client
func initProtocolClient(ctx context.Context, init *clientInit) (err error) {
	if init.c.ProtocolService, err = waasv1.NewProtocolServiceClient(
		ctx,
		init.opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS protocol client: %w", err)
	}
	return
}

// initMpcKeyClient creates a new MPC key service client
func initMpcKeyClient(ctx context.Context, init *clientInit) (err error) {
	if init.c.MpcKeyService, err = waasv1.NewMPCKeyServiceClient(
		ctx,
		init.opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS mpc key client: %w", err)
	}
	return
}

// initBlockchainClient creates a new blockchain service client
func initBlockchainClient(ctx context.Context, init *clientInit) (err error) {
	if init.c.BlockchainService, err = waasv1.NewBlockchainServiceClient(
		ctx,
		init.opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS blockchain client: %w", err)
	}
	return
}

// initMpcWalletClient creates a new MPC wallet service client
func initMpcWalletClient(ctx context.Context, init *clientInit) (err error) {
	if init.c.MpcWalletService, err = waasv1.NewMPCWalletServiceClient(
		ctx,
		init.opts...,
	); err != nil {
		err = fmt.Errorf("unable to init WaaS mpc wallet client: %w", err)
	}
	return
}

// initMpcTransactionClient creates a new MPC transaction service client
func initMpcTransactionClient(ctx context.Context, init *clientInit) (err error) {
	if init.c.MpcTransactionService, err = waasv1.NewMPCTransactionServiceClient(
		ctx,
		init.opts...,
	); err != nil {
		err = fmt.Errorf("unable to init WaaS mpc transaction client: %w", err)
	}
	return
}

// initPoolClient creates a new pool service client
func initPoolClient(ctx context.Context, init *clientInit) (err error) {
	if init.c.PoolService, err = waasv1.NewPoolServiceClient(
		ctx,
		init.opts...,
	); err != nil {
		err = fmt.Errorf("Unable to init WaaS pool client: %w", err)
	}
	return
}

func waasClientDefaults(
	config config.AppConfig,
) (opts []waasClients.WaaSClientOption) {

	apiPrivateKey, err := base64.StdEncoding.DecodeString(config.ApiPrivateKey)
	if err != nil {
		log.Fatalf("Cannot base64 decode private key: %v", err)
	}

	opts = []waasClients.WaaSClientOption{
		waasClients.WithAPIKey(
			&auth.APIKey{
				Name:       config.ApiKeyName,
				PrivateKey: string(apiPrivateKey),
			},
		),
	}
	return
}
