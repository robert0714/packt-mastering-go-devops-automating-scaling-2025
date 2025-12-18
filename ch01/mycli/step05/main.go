package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "A simple CLI application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to the app!")
		},
	}

	// Adding a 'hello' subcommand
	helloCmd := &cobra.Command{
		Use:   "hello",
		Short: "Greet the user",
		Run: func(cmd *cobra.Command, args []string) {
			name := "world"
			if len(args) > 0 {
				name = args[0]
			}
			fmt.Printf("Hello, %s!\n", name)
		},
	}
	rootCmd.AddCommand(helloCmd)

	var greeting string
	helloCmd.Flags().StringVarP(
		&greeting, "greeting", "g", "Hello", "custom greeting")
	helloCmd.Run = func(cmd *cobra.Command, args []string) {
		name := "world"
		if len(args) > 0 {
			name = args[0]
		}
		fmt.Printf("%s, %s!\n", greeting, name)
	}

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
