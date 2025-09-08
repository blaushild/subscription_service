package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	DbHost            string `env:"POSTGRES_HOST"`
	DbPort            string `env:"POSTGRES_PORT"`
	DbName            string `env:"POSTGRES_DB"`
	DbUser            string `env:"POSTGRES_USER"`
	DbPassword        string `env:"POSTGRES_PASSWORD"`
	DbConnectAttempts int    `env:"POSTGRES_CONNECTION_ATTEMPTS"`

	HttpServerPort string `env:"HttpServerPort"`
}

func NewConfig(path ...string) (*Config, error) {
	conf := &Config{}

	// loads from .env file put into ENV for this process
	if err := godotenv.Load(path...); err != nil {
		return nil, err
	}

	// reads from ENV to a struct
	if err := cleanenv.ReadEnv(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
