package main

import (
	"fmt"
	"log"
	"nullableocean-postupashki/src/api/rest"
	"nullableocean-postupashki/src/config"
	_ "nullableocean-postupashki/src/docs"
	"nullableocean-postupashki/src/pkg/hasher"
	"nullableocean-postupashki/src/repository/rabbitmq"
	"nullableocean-postupashki/src/repository/ramstorage"
	"nullableocean-postupashki/src/server"
	"nullableocean-postupashki/src/usecases/service"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title CompileSys
// @version 1.0
// @description This is a compile and execute system.
// @securityDefinitions.apikey APIKeyHeader
// @in header
// @name Authorization
// @description Header expamle: "Authorization: Bearer {token}"
// @BasePath /
func main() {
	appFlags := config.ParseFlags()
	cfg := config.NewAppConfig(appFlags.ConfigPath)

	passHasher := &hasher.BcryptHasher{}

	taskRepo := ramstorage.NewTaskRepository()
	resultRepo := ramstorage.NewResultRepository()
	userRepo := ramstorage.NewUserRepository()
	sessionRepo := ramstorage.NewSessionRepository()

	taskSender, err := rabbitmq.NewRabbitMQTaskSender(cfg.RabbitMQ.GetAmqpUrl(), cfg.RabbitMQ.QueueName)
	if err != nil {
		log.Fatalf("message broker error: %s", err)
	}

	sessionService := service.NewSessionService(sessionRepo)
	userService := service.NewUserService(userRepo, sessionService, passHasher)
	taskService := service.NewTaskService(taskRepo, taskSender)
	resultService := service.NewResultService(resultRepo)

	userHandler := rest.NewUserHandler(userService)
	taskHandler := rest.NewTaskHandler(taskService, resultService, sessionService)
	commitHandler := rest.NewCommitHandler(cfg.Commiter.AccessHeader, cfg.Commiter.AccessToken, resultService)

	router := chi.NewRouter()
	taskHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)
	commitHandler.RegisterRoutes(router)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	server := server.NewServer(cfg.Server.Port, router)

	fmt.Printf("Server listen on http://%s:%s\n...", cfg.Server.Host, cfg.Server.Port)
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
