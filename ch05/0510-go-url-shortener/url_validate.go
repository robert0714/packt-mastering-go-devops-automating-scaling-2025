package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

// Middleware for validating the URL
func validateURLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the full body once
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(
				w, "Unable to read request body",
				http.StatusInternalServerError,
			)
			return
		}
		// Restore the body so it can be read again later
		r.Body = io.NopCloser(bytes.NewReader(data))
		// Check URL in request body
		var urlReq URL
		err = json.Unmarshal(data, &urlReq)
		if err != nil || urlReq.Long == "" {
			http.Error(
				w, "Invalid JSON or empty URL", http.StatusBadRequest,
			)
			return
		}

		// Call httpstatus.io API to validate URL
		if !isValidURL(urlReq.Long) {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		// Call the next handler if URL is valid
		next.ServeHTTP(w, r)
	})
}

// Function to check URL using httpstatus.io
func isValidURL(url string) bool {
	requestBody := map[string]string{
		"requestUrl": url,
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(
		"POST", "https://api.httpstatus.io/v1/status", bytes.NewBuffer(
			body,
		),
	)
	if err != nil {
		log.Println("Error creating request:", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	billingToken := os.Getenv("HTTPSTATUS_TOKEN")
	if billingToken == "" {
		log.Println("HTTPSTATUS_TOKEN environment variable not set")
		return false
	}
	req.Header.Set("X-Billing-Token", billingToken)
	// Send the request to validate the URL
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()
	// We only care about the status code to check URL validity
	return resp.StatusCode == http.StatusOK
}
