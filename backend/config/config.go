package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP         `yaml:"http"`
	PG           `yaml:"pg"`
	Hasher       `yaml:"hasher"`
	TokenManager `yaml:"token_manager"`
	Log          `yaml:"logger"`
}

type HTTP struct {
	Host string `env-required:"true" yaml:"host" env:"HOST"`
	Port string `env-required:"true" yaml:"port" env:"PORT"`
}

type PG struct {
	PoolMax    int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
	URL        string `env-required:"true" yaml:"pg_url"   env:"PG_URL"`
	Schema_Url string `env-required:"true" yaml:"schema_url"   env:"SCHEMA_URL"`
}

type Hasher struct {
	Salt string `env-required:"true" yaml:"salt" env:"PASS_SALT"`
}

type TokenManager struct {
	SigningKey      string        `env-required:"true" yaml:"signing_key" env:"SIGNING_KEY"`
	AccessTokenTTL  time.Duration `env-required:"true" yaml:"atttl" env:"ATTL"`
	RefreshTokenTTL time.Duration `env-required:"true" yaml:"rtttl" env:"RTTL"`
}

type Log struct {
	Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
