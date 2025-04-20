package replication

import "inmemorykvdb/internal/network"

type Master struct {
	server *network.Server
}

func NewMaster() (*Master, error)

func (m *Master) Start()
