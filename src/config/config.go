package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string
	Port string
}

func ReadConfig() *Config {
	envLoad()

	return &Config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}
}

func envLoad() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error load .env file", err)
	}
}
