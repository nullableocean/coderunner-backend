package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	defaultConfigPath = "./src/config/configs/default.yml"
	defaultEnvPath    = "./.env"
)

type AppFlags struct {
	ConfigPath string
}

func ParseFlags() *AppFlags {
	configPath := flag.String("config", defaultConfigPath, "Path to config")

	return &AppFlags{
		ConfigPath: *configPath,
	}
}

type Server struct {
	Host string `yaml:"host" env:"APP_HOST"`
	Port string `yaml:"port" env:"APP_PORT"`
}
type RabbitMQ struct {
	Host      string `yaml:"host" env:"RABBITMQ_HOST"`
	Port      string `yaml:"port" env:"RABBITMQ_PORT"`
	User      string `yaml:"user" env:"RABBITMQ_DEFAULT_USER"`
	Password  string `yaml:"password" env:"RABBITMQ_DEFAULT_PASS"`
	QueueName string `yaml:"queue_name" env:"QUEUE_NAME"`
}

func (cfg *RabbitMQ) GetAmqpUrl() string {
	return fmt.Sprintf("amqp://%s:%s@localhost:5672/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
}

type AppConfig struct {
	RabbitMQ `yaml:"rabbitmq"`
	Server   `yaml:"server"`
}

func NewAppConfig(cfgPath string) *AppConfig {
	loadEnv()

	cfg := &AppConfig{}

	if cfgPath == "" {
		log.Fatal("config path is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist by this path: %s", cfgPath)

	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return cfg
}

func loadEnv() {
	err := godotenv.Load(defaultEnvPath)
	if err != nil {
		log.Fatalf("error load .env file: %s", err)
	}
}
