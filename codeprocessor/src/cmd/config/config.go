package config

import (
	"flag"
	"log"
	"os"
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
	Rabbit   RabbitMqConfig
	Process  ProcessorConfig
	Commiter CommiterConfig
}

type AppFlags struct {
	RunnerDockerfileDir string
}

func ParseFlags() *AppFlags {
	runnerDocDir := flag.String("runnerdoc", "", "Path to config")

	flag.Parse()

	_, err := os.Stat(*runnerDocDir)
	if err != nil {
		log.Fatalln("error path for runner dockerfile directory")
	}

	return &AppFlags{
		RunnerDockerfileDir: *runnerDocDir,
	}
}

func InitConfig() *AppConfig {
	cfg := &AppConfig{}

	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		log.Fatalf("read config error %s", err)
	}

	checkProcessConfig(cfg)

	return cfg
}

func checkProcessConfig(cfg *AppConfig) {
	if cfg.Process.ProcessesLimit == 0 {
		cfg.Process.ProcessesLimit = runtime.NumCPU()
	}
}
