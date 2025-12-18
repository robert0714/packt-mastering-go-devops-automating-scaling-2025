package main

import (
	"fmt"
	"net/http"
	"sync"
)

func fetchURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching:", url, err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Fetched", url, "with status", resp.Status)
}
func main() {
	urls := []string{
		"https://golang.org", "https://go.dev", "https://godoc.org"}
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, &wg)
	}
	wg.Wait()
	fmt.Println("All URLs fetched")
}
