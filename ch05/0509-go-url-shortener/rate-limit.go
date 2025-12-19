package main

import (
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	requests      = make(map[string]int)
	requestsMutex sync.Mutex
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr // fallback if parsing fails
		}
		requestsMutex.Lock()
		requests[ip]++
		count := requests[ip]
		requestsMutex.Unlock()
		if count > 5 {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		time.AfterFunc(time.Minute, func() {
			requestsMutex.Lock()
			requests[ip]--
			requestsMutex.Unlock()
		})
		next.ServeHTTP(w, r)
	})
}
