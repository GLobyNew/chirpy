package main

import (
	"io"
	"net/http"
)

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

func main() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/healthz", handleReadiness)

	serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./app/"))))
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()

}
