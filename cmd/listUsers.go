/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "listUsers",
	Short: "List all the available users/contacts",
	Long:  `List all the available users/contacts.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nList of all the available users/contacts:\n\n")
		listDirectories()
	},
}

func init() {
	rootCmd.AddCommand(listUsersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listUsersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listUsersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listDirectories() {
	storageDir := "storage" // Replace with the actual path to the storage directory

	files, err := os.ReadDir(storageDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fmt.Println("|----------------------|")
	fmt.Println("| Available Users/Contacts |")
	fmt.Println("|----------------------|")

	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("| %-20s |\n", file.Name())
		}
	}
}
