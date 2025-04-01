package config

import (
	"errors"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Engine  *EngineConfig  `yaml:"engine"`
	Network *NetworkConfig `yaml:"network"`
	Logging *LoggingConfig `yaml:"logging"`
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
