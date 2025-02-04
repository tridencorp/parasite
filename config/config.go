package config

import (
	"os"
)

// Simple config that loads and parse json files.
type Config struct {
	data []byte
}

// Create new config and load json file.
func NewConfig(file string) (*Config, error) {
	conf := &Config{}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	conf.data = data
	return conf, nil
}
