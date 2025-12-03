package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filePath = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filePath)))

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Servering file at '%s' on port %s ...\n", filePath, port)
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}
