package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// TODO: DbURL will be set in later stage
type AppConfig struct {
	DbURL string  `env:"DB_URL" envDefault:""`
	Host  string  `env:"HOST" envDefault:"localhost"`
	Port  string  `env:"PORT" envDefault:"8080"`
	Db    *sql.DB `env:"-"`
}

func NewAppConfig() *AppConfig {
	var c AppConfig

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	if err := env.Parse(&c); err != nil {
		log.Fatalf("Couldn't load config file. Error: %s", err.Error())
	}

	fmt.Printf("Config: %+v\n", c)

	return &c
}

func (c *AppConfig) ListnerAddr() string {
	return c.Host + ":" + c.Port
}
