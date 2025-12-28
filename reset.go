package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	cfg.fileserverHits.Store(0)
	resetHits := fmt.Sprintf("Reset Hits: %v\nAll users deleted\n", cfg.fileserverHits.Load())

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.db.DeleteAllUsers(r.Context())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resetHits))
}
