/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// For slog logger check the root.go file in the same directory

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "listUsers",
	Short: "List all the available users/contacts",
	Long:  `List all the available users/contacts.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nList of all the available users/contacts:\n\n")
		slog.Info("Calling the listDirectories function")
		ListDirectories()
		slog.Info("Finshed called the listDirectories function")
	},
}

func init() {
	rootCmd.AddCommand(listUsersCmd)
}

func ListDirectories() {
	storageDir := "storage" // Replace with the actual path to the storage directory

	files, err := os.ReadDir(storageDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fmt.Println("|----------------------|")
	fmt.Println("| Available Users/Contacts |")
	fmt.Println("|----------------------|")
	slog.Info("beginning to list directories for loop")
	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("| %-20s |\n", file.Name())
		}
	}
	slog.Info("ending list directories for loop")
}
