package config_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/fffzlfk/distrikv/config"
)

func createConfig(t *testing.T, contents string) *config.Config {
	f, err := ioutil.TempFile(os.TempDir(), "config.toml")
	if err != nil {
		t.Fatalf("cound not create temp file: %v", err)
	}
	defer f.Close()

	name := f.Name()
	defer os.Remove(name)

	if _, err = f.WriteString(contents); err != nil {
		t.Fatalf("cound not write contents: %v", err)
	}

	cfg, err := config.ParseFile(name)
	if err != nil {
		t.Fatalf("cound not parsefile: %v", err)
	}
	return cfg
}

func TestConfigParse(t *testing.T) {
	contents := `[[shards]]
	name = "Xian"
	index = 0
	address = "localhost:8080"`

	got := createConfig(t, contents)

	want := &config.Config{
		Shards: []config.Shard{
			{
				Name:    "Xian",
				Index:   0,
				Address: "localhost:8080",
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parse failed, want: %#v, get: %#v", want, got)
	}
}

func TestParseShards(t *testing.T) {
	cfg := createConfig(t, `
	[[shards]]
		name = "Beijing"
		index = 0
		address = "localhost:8080"
	[[shards]]
		name = "Shanghai"
		index = 1
		address = "localhost:8081"`)

	got, err := config.ParseShards(cfg.Shards, "Shanghai")
	if err != nil {
		t.Fatalf("could not ParseShards %#v: %v", got, err)
	}

	want := &config.Shards{
		Count: 2,
		Index: 1,
		Addrs: map[int]string{
			0: "localhost:8080",
			1: "localhost:8081",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("parse failed, want: %#v, get: %#v", want, got)
	}
}
