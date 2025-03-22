package main

import (
	"net/http"
	"strconv"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middleWareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(int32(1))
		next.ServeHTTP(w, r)
	})
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	w.Write([]byte("Hits: " + strconv.Itoa(int(hits)) + "\n"))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func main() {
	var apiCfg apiConfig
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/healthz", handleReadiness)
	appHandler := http.StripPrefix("/app/", http.FileServer(http.Dir("./app/")))
	serveMux.Handle("/app/", apiCfg.middleWareMetricsInc(appHandler))
	serveMux.HandleFunc("/metrics", apiCfg.handleMetrics)
	serveMux.HandleFunc("/reset", apiCfg.handleReset)
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()

}
