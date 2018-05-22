package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quipo/statsd"
)

var healthy = true

func main() {
	prefix := "dummy."
	statsdclient := statsd.NewStatsdClient("localhost:8125", prefix)
	err := statsdclient.CreateSocket()
	if nil != err {
		log.Println(err)
		os.Exit(1)
	}
	interval := time.Second * 10
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	defer stats.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		stats.Incr("http.requests.get.root", 1)

		hostname, err := os.Hostname()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write([]byte(hostname))
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		stats.Incr("http.requests.get.healthz", 1)

		if healthy {
			w.Write([]byte("healthy"))
		} else {
			http.Error(w, "unhealthy", http.StatusInternalServerError)
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		stats.Incr("http.requests.post.healthz", 1)

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
