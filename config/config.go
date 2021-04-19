package config

import (
	"bufio"
	"os"

	toml "github.com/pelletier/go-toml"
)

// Shard describes a shard that holds the approprite set of keys
// Each Shard has unique set of keys
type Shard struct {
	Name    string
	Index   int
	Address string
}

// Config describes the sharding config
type Config struct {
	Shards []Shard
}

func ParseFile(configFileName string) (*Config, error) {
	configFile, err := os.Open(configFileName)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := toml.NewDecoder(bufio.NewReader(configFile)).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
