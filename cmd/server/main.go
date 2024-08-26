package main

import(
	"log"
	"github.com/lyteabovenyte/mock_distributed_services/internal/server"
)

func main() {
	srv := server.NewHttpServer(":8080")
	log.Fatal(srv.ListenAndServe())
}