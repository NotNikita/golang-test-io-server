package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Service struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		File     string `yaml:"file"`
		Interval int    `yaml:"interval"`
	} `yaml:"service"`
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig error occured while reading config: %w", err)
	}
	defer f.Close()

	config := &Config{}
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, fmt.Errorf("config.LoadConfig error occured while decoding config: %w", err)
	}

	return config, nil
}
