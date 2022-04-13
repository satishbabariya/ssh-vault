package config

import "github.com/caarlos0/env/v6"

type Config struct {
	VaultSecret string `env:"VAULT_SECRET,required"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	Port        string `env:"PORT" envDefault:"1203"`
}

func NewConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
