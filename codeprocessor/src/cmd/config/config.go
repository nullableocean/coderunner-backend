package config

import (
	"flag"
	"log"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type RabbitMqConfig struct {
	AmqpUrl   string `env:"AMQP_URL"`
	QueueName string `env:"TASKS_QUEUE_NAME"`
}

type ProcessorConfig struct {
	ProcessesLimit int `env:"PROCESSES_LIMIT"`
}

type CommiterConfig struct {
	CommitUrl    string `env:"COMMIT_API_URL"`
	AccessToken  string `env:"COMMIT_API_TOKEN"`
	AccessHeader string `env:"COMMIT_ACCESS_HEADER"`
}

type AppConfig struct {
	RabbitMqConfig
	ProcessorConfig
	CommiterConfig
}

type AppFlags struct {
	RunnerDockerfileDir string
}

func ParseFlags() *AppFlags {
	runnerDocDir := flag.String("runnerdoc", "", "Path to config")

	return &AppFlags{
		RunnerDockerfileDir: *runnerDocDir,
	}
}

func InitConfig() AppConfig {
	cfg := AppConfig{}

	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		log.Fatalf("read config error %s", err)
	}

	if cfg.ProcessesLimit == 0 {
		cfg.ProcessesLimit = runtime.NumCPU()
	}

	return cfg
}
