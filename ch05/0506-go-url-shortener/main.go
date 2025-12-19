package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(LoggingMiddleware) // Apply logging middleware globally
	router.HandleFunc("/", homeHandler).Methods("GET")

	// Creating the URL shortening logic
	router.HandleFunc("/shorten", createShortURLHandler).Methods("POST")
	// Creating the URL shortening logic
	router.HandleFunc("/{short}", getOriginalURLHandler).Methods("GET")

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func homeHandler(response http.ResponseWriter, r *http.Request) {
	response.Write([]byte("Welcome to the URL shortener API!"))
}
