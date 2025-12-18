package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'hello' subcommand")
		return
	}
	command := os.Args[1]
	switch command {
	case "hello":
		name := "world"
		if len(os.Args) > 2 {
			name = os.Args[2]
		}
		fmt.Printf("Hello, %s!\n", name)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
