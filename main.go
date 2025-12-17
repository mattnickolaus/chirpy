package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const filePath = "."

	// var hits atomic.Int32

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetric(http.FileServer(http.Dir(filePath)))))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.numberOfHits)
	mux.HandleFunc("POST /api/reset", apiCfg.resetHits)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Servering file at '%s' on port %s ...\n", filePath, port)
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}
