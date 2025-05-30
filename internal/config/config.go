package config

import (
	"errors"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Engine      *EngineConfig  `yaml:"engine"`
	Network     *NetworkConfig `yaml:"network"`
	Logging     *LoggingConfig `yaml:"logging"`
	WalConfig   *WalConfig     `yaml:"wal"`
	Replication *ReplicaConfig `yaml:"replication"`
}

func NewConfig(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, errors.New("reader is nil")
	}

	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, errors.New("problems with reading file")
	}

	var config Config

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, errors.New("problems with parsing file")
	}

	return &config, nil
}

type EngineConfig struct {
	EngineType string `yaml:"type"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
	IsSync         bool          `yaml:"is_sync"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type WalConfig struct {
	BatchSize      int           `yaml:"flushing_batch_size"`
	BatchTimeout   time.Duration `yaml:"flushing_batch_timeout"`
	MaxSegmentSize string        `yaml:"max_segment_size"`
	DataDirectory  string        `yaml:"data_directory"`
	FileName       string        `yaml:"file_name"`
}

type ReplicaConfig struct {
	ReplicaType   string        `yaml:"replica_type"`
	MasterAddress string        `yaml:"master_address"`
	SyncInterval  time.Duration `yaml:"sync_interval"`
}

type ClientConfig struct {
	Address        string
	MaxMessageSize int
	Timeout        time.Duration
}
