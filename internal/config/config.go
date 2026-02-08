package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for repocheck
type Config struct {
	IgnorePaths []string `yaml:"ignore_paths"`
	StrictMode  bool     `yaml:"strict_mode"`
}

// LoadConfig reads the .repocheck.yaml file if it exists
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = ".repocheck.yaml"
	}

	cfg := &Config{
		IgnorePaths: []string{}, // defaults
		StrictMode:  false,
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil // return default config if no file found
	}
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
