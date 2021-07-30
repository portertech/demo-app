package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quipo/statsd"
)

var healthy = true

var db_host = os.Getenv("DB_HOST")
var db_port = os.Getenv("DB_PORT")
var db_user = os.Getenv("DB_USER")
var db_pass = os.Getenv("DB_PASS")
var db_name = os.Getenv("DB_NAME")

func main() {
	prefix := "demo-app."
	statsdclient := statsd.NewStatsdClient("127.0.0.1:8125", prefix)
	err := statsdclient.CreateSocket()
	if nil != err {
		log.Println(err)
		os.Exit(1)
	}
	interval := time.Second * 5
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	defer stats.Close()

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_pass, db_name)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

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

		err = db.Ping()
		if err != nil {
			healthy = false
			log.Println(err)
		}

		if healthy {
			w.Write([]byte("healthy"))
		} else {
			http.Error(w, "unhealthy", http.StatusInternalServerError)
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		stats.Incr("http.requests.post.healthz", 1)

		err = db.Ping()
		if err != nil {
			healthy = false
			log.Println(err)
		} else {
			healthy = !healthy
		}
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
