package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://golang.org")
	if err != nil {
		fmt.Println("Request timed out:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
}
