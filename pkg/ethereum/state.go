// Package ethereum contains all the functions related to Ethereum State and Synchronization
// with Layer 1
package ethereum

import (
	"context"
	"github.com/NethermindEth/juno/internal/log"
	"github.com/NethermindEth/juno/pkg/db"
	"github.com/NethermindEth/juno/pkg/feeder_gateway"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/core/types"
)

// Synchronizer represents the base struct for Ethereum Synchronization
type Synchronizer struct {
	ethereumClient      *ethclient.Client
	feederGatewayClient *feeder_gateway.Client
	db                  *db.Databaser
}

// NewSynchronizer creates a new Synchronizer
func NewSynchronizer(ethereumNode, feederGateway string, db *db.Databaser) *Synchronizer {
	client, err := ethclient.Dial(ethereumNode)
	if err != nil {
		log.Default.With("Error", err).Fatal("Unable to connect to Ethereum Client")
	}
	feeder := feeder_gateway.NewClient(feederGateway, "/feeder_gateway", nil)
	return &Synchronizer{
		ethereumClient:      client,
		feederGatewayClient: feeder,
		db:                  db,
	}
}

// UpdateStateRoot keeps updating the Ethereum State Root as a process
func (s Synchronizer) UpdateStateRoot() error {
	headers := make(chan *types.Header)
	sub, err := s.ethereumClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Default.With("Error", err).Fatalf("Unable to subscribe to block headers")
		return err
	}

	for {
		select {
		case err := <-sub.Err():
			log.Default.Fatal(err)
			return err
		case header := <-headers:
			log.Default.With("stateRoot", header.Root).Debug("State root retrieved from L1")
			// TODO store ethereum state
		}
	}
}

// Close closes the client for the Layer 1 Ethereum node
func (s Synchronizer) Close(ctx context.Context) {
	// notest
	log.Default.Info("Closing Layer 1 Synchronizer")
	select {
	case <-ctx.Done():
		s.ethereumClient.Close()
	default:
	}
}