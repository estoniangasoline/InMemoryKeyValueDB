package initialization

import (
	"errors"
	"inmemorykvdb/internal/config"
	"inmemorykvdb/internal/database/storage/replication"
	"inmemorykvdb/internal/network"

	"go.uber.org/zap"
)

const (
	slave  = "slave"
	master = "master"
)

func createReplica(logger *zap.Logger, replCnfg *config.ReplicaConfig, walCnfg *config.WalConfig) (replica, error) {
	if logger == nil {
		return nil, errors.New("could not create replica without logger")
	}

	if replCnfg == nil {
		return nil, nil
	}

	switch replCnfg.ReplicaType {
	case slave:
		client, err := network.NewClient(replCnfg.MasterAddress)

		if err != nil {
			return nil, errors.New("could not create client for slave")
		}

		return replication.NewSlave(client, logger, replication.WithInterval(replCnfg.SyncInterval),
			replication.WithDirectorySlave(walCnfg.DataDirectory))
	case master:
		server, err := network.NewServer(replCnfg.MasterAddress, logger)

		if err != nil {
			return nil, errors.New("could not create server for master")
		}

		return replication.NewMaster(server, logger,
			replication.WithDirectoryMaster(walCnfg.DataDirectory))
	}

	return nil, errors.New("unknown replica type")
}
