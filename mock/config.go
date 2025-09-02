package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Port string `yaml:"port"`
}

func LoadConfig(path string) (*Configuration, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cfg Configuration

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
