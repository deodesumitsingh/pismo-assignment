package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

// TODO: DbURL will be set in later stage
type AppConfig struct {
	DbURL string `env:"DB_URL" envDefault:""`
	Host  string `env:"HOST" envDefault:"localhost"`
	Port  string `env:"PORT" envDefault:"8080"`
}

func NewAppConfig() *AppConfig {
	var c AppConfig

	if err := env.Parse(&c); err != nil {
		log.Fatalf("Couldn't load config file. Error: %s", err.Error())
	}

	return &c
}

func (c *AppConfig) ListnerAddr() string {
	return c.Host + ":" + c.Port
}
