package main

import (
	"fmt"
	"log"
	"nullableocean-postupashki/src/api/rest/compilehandler"
	_ "nullableocean-postupashki/src/docs"
	"nullableocean-postupashki/src/server"
	"nullableocean-postupashki/src/service/compile"
	"nullableocean-postupashki/src/storage/ramstorage"

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
	db := ramstorage.NewRamStorage()
	cmplService := compile.NewCompileService(db)
	cmplHandler := compilehandler.NewCompileHandler(cmplService)
	router := chi.NewRouter()
	cmplHandler.RegisterRoutes(router)

	router.Get("/swagger/*", httpSwagger.WrapHandler)
	server := server.NewServer(PORT, router)

	fmt.Printf("Server listen on http://%s:%s\n...", HOST, PORT)

	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
