package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)
	resetHits := fmt.Sprintf("Reset Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(resetHits))
}
