package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Listening...\n")
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}
