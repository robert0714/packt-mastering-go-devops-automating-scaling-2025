package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Create router and register routes
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/shorten", createShortURLHandler).Methods("POST")
	r.HandleFunc("/{short}", getOriginalURLHandler).Methods("GET")
	// Attach the validateURLMiddleware to the /shorten route
	r.Handle("/shorten", validateURLMiddleware(http.HandlerFunc(createShortURLHandler)))
	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())
	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(response http.ResponseWriter, r *http.Request) {
	response.Write([]byte("Welcome to the URL shortener API!"))
}

func getStatsHandler(response http.ResponseWriter, r *http.Request) {
	response.Write([]byte("This is getStatsHandler!"))
}
