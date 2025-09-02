package main

import (
	"sync"

	"os"

	"gopkg.in/yaml.v3"
)

type ConfigurationStore struct {
	mu   sync.RWMutex
	conf *Configuration
}

type Configuration struct {
	APIKey   string            `yaml:"apiKey"`
	Services map[string]string `yaml:"services"`
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

func NewConfigStore(c *Configuration) *ConfigurationStore {
	return &ConfigurationStore{conf: c}
}

func (cs *ConfigurationStore) Get() *Configuration {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.conf
}

func (cs *ConfigurationStore) Update(c *Configuration) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.conf = c
}
