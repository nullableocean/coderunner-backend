package main

import (
	"codeproccesor/src/cmd/config"
	"codeproccesor/src/usecases/service/coderunner"
	"codeproccesor/src/usecases/service/commiter/httpcommiter"
	"codeproccesor/src/usecases/service/consumer"
	"codeproccesor/src/usecases/service/processor"
	"codeproccesor/src/usecases/service/processpool"
	"fmt"
	"sync"

	"context"
	"log"
	"os/signal"
	"syscall"

	dockerClient "github.com/docker/docker/client"
)

var ()

func main() {
	flags := config.ParseFlags()
	cfg := config.InitConfig()

	logger := log.Default()

	docClient, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	runner, err := coderunner.NewCodeRunner(docClient, flags.RunnerDockerfileDir)
	if err != nil {
		log.Fatal(err)
	}

	commiter := httpcommiter.NewHttpCommiter(cfg.CommitUrl, cfg.CommiterConfig.AccessToken, cfg.CommiterConfig.AccessHeader)
	codeProcessor := processor.NewCodeProcessor(commiter, runner)
	processPool := processpool.NewProcessPool(cfg.ProcessesLimit, codeProcessor, logger)
	taskConsumer := consumer.NewTaskConsumer(processPool, logger)

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		<-ctx.Done()

		fmt.Println("stopping service...")

		taskConsumer.Stop()
		processPool.Stop()

		fmt.Println("service stoped...")
		wg.Done()
	}()

	err = taskConsumer.Consume(cfg.AmqpUrl, cfg.QueueName)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
