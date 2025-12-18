package main

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	bar := progressbar.NewOptions(100,
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(20),
	)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(50 * time.Millisecond)
	}
}
