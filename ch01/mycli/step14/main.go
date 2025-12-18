package main

import (
	"fmt"
	"time"
)

func main() {
	rate := time.Second / 5 // 5 requests per second
	ticker := time.NewTicker(rate)
	defer ticker.Stop()
	for i := 0; i < 20; i++ {
		<-ticker.C
		fmt.Println("Request", i+1)
	}
}
