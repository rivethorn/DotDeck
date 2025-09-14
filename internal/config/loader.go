// Package config holds the functionality for reading the TOML config files.
package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

// Config holds a map of source and destination files to process
type Config struct {
	Files map[string]string `toml:"files"`
}

// Load loads the .toml config file and parses it
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
