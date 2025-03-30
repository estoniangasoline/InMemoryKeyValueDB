package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/network"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	byteMultiply     = 1
	kilobyteMultiply = 1024
	megabyteMultiply = 1048576

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
			network.WithServerTimeout(defaultTimeout),
			network.WithSync(defaultIsSync))
	}

	address := cnfg.Address

	if address == "" {
		address = defaultAddress
	}

	maxConnections := cnfg.MaxConnections

	if maxConnections == 0 {
		maxConnections = defaultMaxConnections
	}

	maxMsgSize := defaultMaxMessageSize
	multiply := defaultMultiply

	if len(cnfg.MaxMessageSize) > 2 {

		unparsedSize := cnfg.MaxMessageSize[len(cnfg.MaxMessageSize)-2:]
		probablyMultiply, err := parseMessageSize(unparsedSize)

		if err == nil {
			multiply = probablyMultiply

			probablyMessageSize, err := strconv.Atoi(cnfg.MaxMessageSize[:len(cnfg.MaxMessageSize)-2])

			if err == nil {
				maxMsgSize = probablyMessageSize * multiply
			}
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

func parseMessageSize(unparsedSize string) (int, error) {

	if unparsedSize[0] <= '9' && '0' <= unparsedSize[0] && unparsedSize[1] == 'B' {
		return byteMultiply, nil
	}

	switch unparsedSize {
	case "KB":
		return kilobyteMultiply, nil
	case "MB":
		return megabyteMultiply, nil
	}

	return -1, errors.New("unknown or forbidden buffer size")
}
