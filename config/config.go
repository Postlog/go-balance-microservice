package config

import (
	"encoding/json"
	"github.com/caarlos0/env/v6"
	"os"
)

// Config holds application configuration values
type Config struct {
	Port string `env:"MICROSERVICE_PORT,required"`
	DSN  string `env:"MICROSERVICE_DSN,required"`
	ExchangeRatesAPIKey string `env:"MICROSERVICE_ER_API_KEY,required"`
	ApiRequestTimeout   int64  `json:"apiRequestTimeout"`
	BaseCurrency        string `json:"baseCurrency"`
	Logger              json.RawMessage
}

// Load loads configuration values from environment and JSON file, located on the provided path
func Load(path string) (*Config, error) {
	cfg := Config{}
	if err := cfg.loadEnv(); err != nil {
		return nil, err
	}
	if err := cfg.loadJSON(path); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) loadEnv() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) loadJSON(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytes, cfg); err != nil {
		return err
	}
	return nil
}
