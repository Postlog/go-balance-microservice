package config

import (
	"encoding/json"
	"github.com/caarlos0/env/v6"
	"os"
)

type Config struct {
	Address             string `env:"BALANCE_ADDRESS,required"`
	DatabaseAddress     string `env:"BALANCE_DATABASE,required"`
	ExchangeRatesAPIKey string `env:"BALANCE_ER_API,required"`
	IsDev               bool   `env:"BALANCE_IS_DEV" envDefault:"false"`
	ApiRequestTimeout   int64  `json:"apiRequestTimeout"`
	BaseCurrency        string `json:"baseCurrency"`
	Logger              json.RawMessage
}

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
