package main

import (
	"log"

	"github.com/CepstrumLabs/proglog/internal/server"
)

func main() {
	srv := server.NewHttpServer(":8080")
	log.Println("Starting server at localhost:8080 ... ")
	log.Fatal(srv.ListenAndServe())
}
