/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/cloudflare/circl/sign/dilithium"
	"github.com/spf13/cobra"
)

// userIDCmd represents the userID command
var userIDCmd = &cobra.Command{
	Use:   "userID",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter the user ID: ")
		var userID string
		fmt.Scanln(&userID)

		//check if the user exists
		checkUser(userID)

		// Prompt the user if they want to send the message
		fmt.Print("Do you want to send the message? (y/n): ")
		var sendChoice string
		fmt.Scanln(&sendChoice)

		if sendChoice == "y" {
			fmt.Print("Enter the message: ")
			var message string
			fmt.Scanln(&message)
			// Send the message
			SendMsg(userID, message)
		} else {
			fmt.Println("Message not sent.")
		}

	},
}

func init() {
	rootCmd.AddCommand(userIDCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userIDCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userIDCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkUser(userID string) {
	if _, err := os.Stat("storage/" + userID); os.IsNotExist(err) {
		fmt.Println("User" + userID + " does not exist")
		err := os.Mkdir("storage/"+userID, 0755)
		if err != nil {
			fmt.Println("Failed to create User:"+userID, err)
			return
		}
		fmt.Println("User" + userID + " created successfully")
	} else {
		fmt.Println("User " + userID + " exists")
	}
	checkKeys(userID)
}

func checkKeys(userID string) {

	// Check if the Dilithium5 key pair already exists
	if _, err := os.Stat("storage/" + userID + "/keys/sign/privateKey"); !os.IsNotExist(err) {
		fmt.Println("Key pair already exists for " + userID)
		return
	} else { // else create a new key pair
		modename := "Dilithium5"

		// Generate Dilithium5 key pair
		mode := dilithium.ModeByName(modename)

		publicKey, privateKey, err := mode.GenerateKey(nil)
		if err != nil {
			fmt.Println("Failed to generate Key pair:", err)
			return
		}

		// Create the keys directory if it doesn't exist
		if _, err := os.Stat("storage/" + userID + "/keys/sign"); os.IsNotExist(err) {
			err := os.Mkdir("storage/"+userID+"/keys/sign", 0755)
			if err != nil {
				fmt.Println("Failed to create keys directory for User:"+userID, err)
				return
			}
		}

		// Save the key pair to files
		err = os.WriteFile("storage/"+userID+"/keys/sign/privateKey", privateKey.Bytes(), 0644)
		if err != nil {
			fmt.Println("Failed to save private key:", err)
			return
		}
		err = os.WriteFile("storage/"+userID+"/keys/sign/publicKey", publicKey.Bytes(), 0644)
		if err != nil {
			fmt.Println("Failed to save public key:", err)
			return
		}

		fmt.Println("Key pair created successfully for " + userID)
	}
}
