package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Use(RateLimitMiddleware)
	router.Use(LoggingMiddleware) // Apply logging middleware globally
	// router.Use(AuthenticationMiddleware) // Apply authentication middleware

	// Creating the URL shortening logic
	router.Handle(
		"/shorten",
		validateURLMiddleware(http.HandlerFunc(createShortURLHandler)),
	).Methods("POST")

	// Create a subrouter for /shorten routes
	shortenSubrouter := router.PathPrefix("/shorten").Subrouter()
	shortenSubrouter.Use(validateURLMiddleware)
	shortenSubrouter.HandleFunc("", createShortURLHandler).Methods("POST")

	// We can still define other routes on the main router
	router.HandleFunc("/info", getStatsHandler).Methods("GET")

	router.HandleFunc("/", homeHandler).Methods("GET")

	// Creating the URL shortening logic
	//router.HandleFunc("/shorten", createShortURLHandler).Methods("POST")

	// Creating the URL shortening logic
	router.HandleFunc("/{short}", getOriginalURLHandler).Methods("GET")

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func homeHandler(response http.ResponseWriter, r *http.Request) {
	response.Write([]byte("Welcome to the URL shortener API!"))
}

func getStatsHandler(response http.ResponseWriter, r *http.Request) {
	response.Write([]byte("This is getStatsHandler!"))
}
