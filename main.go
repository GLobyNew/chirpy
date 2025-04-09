package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/GLobyNew/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("error openning db")
	}
	dbQueries := database.New(db)
	platform := os.Getenv("PLATFORM")
	cfg := apiConfig{
		db:       dbQueries,
		platform: platform,
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /admin/healthz", handleReadiness)
	appHandler := http.StripPrefix("/app/", http.FileServer(http.Dir("./app/")))
	serveMux.Handle("/app/", cfg.middleWareMetricsInc(appHandler))
	serveMux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	serveMux.HandleFunc("POST /admin/reset", cfg.handleReset)
	serveMux.HandleFunc("POST /api/users", cfg.handleUser)
	serveMux.HandleFunc("POST /api/login", cfg.handleLogin)
	serveMux.HandleFunc("POST /api/chirps", cfg.handleChirpCreation)
	serveMux.HandleFunc("GET /api/chirps", cfg.handleGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handleGetChirp)
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()

}
