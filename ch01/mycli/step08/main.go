package main

import (
	"fmt"

	"github.com/fatih/color"
)

func main() {
	success := color.New(color.FgGreen).SprintFunc()
	warning := color.New(color.FgYellow).SprintFunc()
	error := color.New(color.FgRed).SprintFunc()
	fmt.Println(success("Operation completed successfully!"))
	fmt.Println(warning("This is a warning message."))
	fmt.Println(error("An error occurred during the operation."))
}
