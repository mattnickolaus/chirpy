package main

import (
	"fmt"
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
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /metrics", apiCfg.numberOfHits)
	mux.HandleFunc("POST /reset", apiCfg.resetHits)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Servering file at '%s' on port %s ...\n", filePath, port)
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) numberOfHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	numberHitsMessage := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(numberHitsMessage))
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)
	resetHits := fmt.Sprintf("Reset Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(resetHits))
}

func (cfg *apiConfig) middlewareMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
