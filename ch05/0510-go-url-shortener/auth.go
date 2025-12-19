package main

import (
	"net/http"
)

// A simple user-role mapping with API keys
var apiKeys = map[string]string{
	"key-admin-123": "admin",
	"key-user-456":  "user",
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		role, ok := apiKeys[apiKey]
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// Store the role in the request context (could be used for RBAC
		// later)
		r.Header.Set("X-User-Role", role)
		next.ServeHTTP(w, r)
	})
}
