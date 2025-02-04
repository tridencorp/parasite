package config

import (
	"encoding/json"
	"os"
)

// Simple config that loads and parse json files.
type Config struct {
	json map[string]json.RawMessage
}

// Create new config based on json file.
func NewConfig(file string) (*Config, error) {
	conf := &Config{}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &conf.json)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

// Unmarshal given configuration key and put it's content to dst.
func (c *Config) Parse(key string, dst any) error {
	return json.Unmarshal(c.json[key], dst)
}
