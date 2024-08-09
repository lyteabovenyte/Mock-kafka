package main

import(
	"log"
	"github.com/lyteabovenyte/distributed_services_with_go/internal/server"
)

func main() {
	srv := server.NewHttpServer(":8080")
	log.Fatal(srv.ListenAndServe())
}