package database

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type DatabaseConfig struct {
	Driver         string `yaml:"driver"`
	DataSourceName string `yaml:"dataSourceName"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

func LoadConfig(filePath string) (*Config, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	return &config, nil
}
