package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	MonitoredAddress string `env:"MONITORED_ADDRESS"`
	JSONRPCServer    string `env:"JSONRPC_SERVER"`

	DBHost     string `env:"DB_HOST"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBPort     int    `env:"DB_PORT"`
	DBDialect  string `env:"DB_DRIVER"`
}

var config = Config{}

func Init() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	err = env.Parse(&config)
	if err != nil {
		panic(err)
	}

}
func V() Config {
	return config
}
