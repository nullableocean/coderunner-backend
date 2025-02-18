package main

import (
	"fmt"
	"log"
	"nullableocean-postupashki/src/api"
	"nullableocean-postupashki/src/server"
	"nullableocean-postupashki/src/service"
	"nullableocean-postupashki/src/storage"

	_ "nullableocean-postupashki/src/docs"

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
	db := storage.NewRamStorage()
	ts := service.NewCompileManager(db)
	th := api.NewTaskHandler(ts)
	router := api.NewChiRouter(th)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", HOST, PORT)),
	))
	server := server.NewServer(PORT, router)

	fmt.Printf("Server listen on http://%s:%s\n...", HOST, PORT)

	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
