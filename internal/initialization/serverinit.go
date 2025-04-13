package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/network"
	"inmemorykvdb/pkg/parsing"
	"time"

	"go.uber.org/zap"
)

const (
	defaultAddress        = ":8080"
	defaultMaxConnections = 100
	defaultMaxMessageSize = 1024
	defaultMultiply       = 4
	defaultTimeout        = 0 * time.Second
	defaultIsSync         = false
)

func createServer(cnfg *config.NetworkConfig, logger *zap.Logger) (*network.Server, error) {

	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	if cnfg == nil {
		return network.NewServer(defaultAddress, logger,
			network.WithServerMaxBufferSize(defaultMaxMessageSize),
			network.WithServerMaxConnections(defaultMaxConnections),
			network.WithServerTimeout(defaultTimeout))
	}

	address := cnfg.Address

	if address == "" {
		address = defaultAddress
	}

	maxConnections := cnfg.MaxConnections

	if maxConnections == 0 {
		maxConnections = defaultMaxConnections
	}

	maxMsgSize := defaultMaxMessageSize * defaultMultiply

	if len(cnfg.MaxMessageSize) > 2 {

		probablyMaxMsgSize, err := parsing.ParseSize(cnfg.MaxMessageSize)

		if err != nil {
			maxMsgSize = probablyMaxMsgSize
		}
	}

	timeOut := defaultTimeout

	if cnfg.IdleTimeout != 0 {
		timeOut = cnfg.IdleTimeout
	}

	server, err := network.NewServer(cnfg.Address, logger,
		network.WithServerMaxBufferSize(maxMsgSize),
		network.WithServerMaxConnections(maxConnections),
		network.WithServerTimeout(timeOut))

	return server, err
}
