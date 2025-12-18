package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mycli",
		Short: "A simple CLI application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to the app!")
		},
	}

	// Adding a 'cloud' subcommand with hierarchical subcommands
	cloudCmd := &cobra.Command{
		Use:   "cloud",
		Short: "Manage cloud resources",
	}
	rootCmd.AddCommand(cloudCmd)

	// cloud create subcommand
	var createName string
	cloudCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new cloud resource",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Creating cloud resource with name: %s\n", createName)
		},
	}
	cloudCreateCmd.Flags().StringVarP(&createName, "name", "n", "", "Name of the resource to create")
	cloudCmd.AddCommand(cloudCreateCmd)

	// cloud delete subcommand
	var deleteID string
	cloudDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a cloud resource",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Deleting cloud resource with id: %s\n", deleteID)
		},
	}
	cloudDeleteCmd.Flags().StringVarP(&deleteID, "id", "i", "", "ID of the resource to delete")
	cloudCmd.AddCommand(cloudDeleteCmd)

	// cloud list subcommand
	cloudListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all cloud resources",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing all cloud resources...")
		},
	}
	cloudCmd.AddCommand(cloudListCmd)

	// cloud update subcommand
	var updateID string
	var updateName string
	cloudUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a cloud resource",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Updating cloud resource with id: %s, new name: %s\n", updateID, updateName)
		},
	}
	cloudUpdateCmd.Flags().StringVarP(&updateID, "id", "i", "", "ID of the resource to update")
	cloudUpdateCmd.Flags().StringVarP(&updateName, "name", "n", "", "New name for the resource")
	cloudCmd.AddCommand(cloudUpdateCmd)

	// Adding a 'deploy' subcommand
	var deployFile string
	var deployTimeout int
	var deployVerbose bool
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy application using configuration",
		Example: `  mycli deploy -f config.yaml
  mycli deploy --timeout 120`,
		Run: func(cmd *cobra.Command, args []string) {
			if deployVerbose {
				fmt.Println("Verbose mode enabled")
			}
			fmt.Printf("Deploying with file: %s, timeout: %d seconds\n", deployFile, deployTimeout)
		},
	}
	deployCmd.Flags().StringVarP(&deployFile, "file", "f", "", "Path to deployment configuration file")
	deployCmd.Flags().IntVarP(&deployTimeout, "timeout", "t", 60, "Maximum wait time")
	deployCmd.Flags().BoolVarP(&deployVerbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.AddCommand(deployCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
