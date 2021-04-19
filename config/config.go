package config

import (
	"bufio"
	"errors"
	"fmt"
	"hash/fnv"
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

// ParseFile loads config from file
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

type Shards struct {
	Count int
	Index int
	Addrs map[int]string
}

// ParseShards provides Shards info from list of shards
func ParseShards(shards []Shard, curShardName string) (*Shards, error) {
	count := len(shards)
	index := -1
	addrs := make(map[int]string)

	for _, v := range shards {
		if _, has := addrs[v.Index]; has {
			return nil, errors.New("duplicated shard index")
		}
		addrs[v.Index] = v.Address
		if v.Name == curShardName {
			index = v.Index
		}
	}

	for i := 0; i < count; i++ {
		if _, has := addrs[i]; !has {
			return nil, fmt.Errorf("shard %d was not found", i)
		}
	}

	if index == -1 {
		return nil, fmt.Errorf("shard %q was not found", curShardName)
	}

	return &Shards{
		Count: count,
		Index: index,
		Addrs: addrs,
	}, nil
}

func (s *Shards) GetIndex(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}
