package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/mattnickolaus/chirpy/internal/database"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func main() {
	godotenv.Load()

	const port = "8080"
	const filePath = "."

	dbURL := os.Getenv("DB_URL")
	platformType := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening databse")
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platformType,
		secret:         secret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetric(http.FileServer(http.Dir(filePath)))))

	mux.HandleFunc("GET /api/healthz", healthHandler)

	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("POST /api/login", apiCfg.login)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeToChirpyRed)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)

	mux.HandleFunc("POST /api/refresh", apiCfg.refreshAccessToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshToken)

	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)
	mux.HandleFunc("GET /admin/metrics", apiCfg.numberOfHits)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Servering file at '%s' on port %s ...\n", filePath, port)
	log.Fatal(http.ListenAndServe(server.Addr, server.Handler))
}
