package main

import (
	"fmt"
	"net/http"
	"time"
)

func fetchWithRetry(url string, attempts int, delay time.Duration) error {
	for i := 0; i < attempts; i++ {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			fmt.Println("Fetched", url, "on attempt", i+1)
			return nil
		}
		fmt.Println("Attempt", i+1, "failed; retrying in", delay)
		time.Sleep(delay)
	}
	return fmt.Errorf(
		"failed to fetch %s after %d attempts", url, attempts)
}
func main() {
	err := fetchWithRetry("https://golang.org", 3, 2*time.Second)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Successfully fetched URL")
	}
}
