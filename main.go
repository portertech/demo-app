package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var healthy = true

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("/www")))
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if healthy {
			w.Write([]byte("healthy"))
		} else {
			http.Error(w, "unhealthy", http.StatusInternalServerError)
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthy = !healthy
	}).Methods(http.MethodPost)

	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
