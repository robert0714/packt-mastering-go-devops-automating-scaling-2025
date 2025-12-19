package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// URL represents the structure of a URL mapping.
// To represent both the input and output of our URL shortener service, we define a simple struct.
// This struct holds the original long URL and the generated short URL:
type URL struct {
	Short string `json:"short_url"`
	Long  string `json:"long_url"`
}

// We use a map to store short URLs and their corresponding long URLs.
// Since this map could be accessed by multiple users at the same time,
// we also define a mutex to make those accesses safe:
var (
	urlStore = make(map[string]string)
	mutex    = &sync.Mutex{}
)

// This function generates a short, random string.
// This will serve as the unique code for each shortened URL.
// It picks six random characters from a predefined set of letters and numbers:
func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	short := make([]byte, 6)
	for i := range short {
		short[i] = letters[rand.Intn(len(letters))]
	}
	return string(short)
}

// The following function handles requests to shorten a long URL.
// It decodes the input JSON, generates a short code, stores the mapping safely,
// and responds with the full shortened URL in JSON format:
func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var request URL
	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate URL input
	if request.Long == "" {
		respondWithError(w, http.StatusBadRequest, "Missing long_url field")
		return
	}

	// Generate short URL and store it
	short := generateShortURL()
	mutex.Lock()
	urlStore[short] = request.Long
	mutex.Unlock()

	// Response with JSON
	response := URL{
		Short: "http://localhost:8080/" + short,
		Long:  request.Long,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// The getOriginalURLHandler function handles requests to a short URL.
// It extracts the short code from the URL path, looks it up in the map, and, if found,
// redirects the user to the original long URL. If not found, it returns an error.
func getOriginalURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	short := vars["short"]
	mutex.Lock()
	longURL, exists := urlStore[short]
	mutex.Unlock()
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
