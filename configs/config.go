package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	Dsn string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
	}
}
