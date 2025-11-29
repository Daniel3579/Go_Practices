package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN     string
	BcryptCost int
	Addr       string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func Load() Config {
	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if cost < 10 || err != nil {
		cost = 12
	}

	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return Config{
		DB_DSN:     os.Getenv("DB_DSN"),
		BcryptCost: cost,
		Addr:       addr,
	}
}
