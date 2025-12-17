package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) numberOfHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	numberHitsMessage := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(numberHitsMessage))
}
