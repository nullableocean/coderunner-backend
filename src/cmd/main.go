package main

import (
	"fmt"
	"log"
	"nullableocean-postupashki/src/api/rest"
	"nullableocean-postupashki/src/config"
	_ "nullableocean-postupashki/src/docs"
	"nullableocean-postupashki/src/pkg/hasher"
	"nullableocean-postupashki/src/repository/ramstorage"
	"nullableocean-postupashki/src/server"
	"nullableocean-postupashki/src/usecases/service"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title CompileSys
// @version 1.0
// @description This is a compile and execute system.
// @BasePath /
func main() {
	cnf := config.ReadConfig()

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

	server := server.NewServer(cnf.Port, router)

	fmt.Printf("Server listen on http://%s:%s\n...", cnf.Host, cnf.Port)
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
