package main

import (
	"fmt"
	"log"
	"nullableocean-postupashki/src/api/rest"
	_ "nullableocean-postupashki/src/docs"
	"nullableocean-postupashki/src/pkg/hasher"
	"nullableocean-postupashki/src/repository/ramstorage"
	"nullableocean-postupashki/src/server"
	"nullableocean-postupashki/src/usecases/service"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	PORT = "8080"
	HOST = "127.0.0.1"
)

// @title CompileSys
// @version 1.0
// @description This is a compile and execute system.

// @host 127.0.0.1:8080
// @BasePath /
func main() {
	passHasher := &hasher.BcryptHasher{}

	taskRepo := ramstorage.NewTaskRepository()
	userRepo := ramstorage.NewUserRepository()
	sessionRepo := ramstorage.NewSessionRepository()

	sessionService := service.NewSessionService(sessionRepo)
	userService := service.NewUserService(userRepo, sessionService, passHasher)
	taskService := service.NewTaskService(taskRepo)

	userHandler := rest.NewUserHandler(userService)
	taskHandler := rest.NewTaskHandler(taskService, sessionService)

	router := chi.NewRouter()
	taskHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	server := server.NewServer(PORT, router)

	fmt.Printf("Server listen on http://%s:%s\n...", HOST, PORT)
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
