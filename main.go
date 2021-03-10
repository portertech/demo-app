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
	prefix := "demo-app."
	statsdclient := statsd.NewStatsdClient("127.0.0.1:8125", prefix)
	err := statsdclient.CreateSocket()
	if nil != err {
		log.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		statsdclient.Incr("http.requests.get.root", 1)

		hostname, err := os.Hostname()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write([]byte(hostname))
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		statsdclient.Incr("http.requests.get.healthz", 1)

		if healthy {
			w.Write([]byte("healthy"))
		} else {
			http.Error(w, "unhealthy", http.StatusInternalServerError)
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		statsdclient.Incr("http.requests.post.healthz", 1)

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
