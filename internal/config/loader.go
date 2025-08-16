package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Files map[string]string `toml:"files"`
}

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
