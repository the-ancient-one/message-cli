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

// sendMsgCmd represents the sendMsg command
var sendMsgCmd = &cobra.Command{
	Use:   "sendMsg",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var userID string
		var message string

		fmt.Print("Enter user ID: ")
		fmt.Scanln(&userID)

		fmt.Print("Enter message: ")
		fmt.Scanln(&message)

		SendMsg(userID, message)
	},
}

func init() {
	rootCmd.AddCommand(sendMsgCmd)

}

func SendMsg(userID string, message string) {
	fmt.Println("Sending message to", userID)
	fmt.Println("Message:", message)

	modename := "Dilithium5"

	mode := dilithium.ModeByName(modename)

	msg := []byte(message)

	if _, err := os.Stat("storage/" + userID + "/keys/privateKey"); !os.IsNotExist(err) {

		pkFile := "storage/" + userID + "/keys/privateKey"

		pubFile := "storage/" + userID + "/keys/publicKey"

		privateKeyBytes, err := os.ReadFile(pkFile)
		if err != nil {
			fmt.Println("Failed to read the private key file:", err)
			return
		}

		publicKeyBytes, err := os.ReadFile(pubFile)
		if err != nil {
			fmt.Println("Failed to read the Public key file:", err)
			return
		}

		privateKey := mode.PrivateKeyFromBytes(privateKeyBytes)

		publiceKey := mode.PublicKeyFromBytes(publicKeyBytes)

		signature := mode.Sign(privateKey, msg)

		fmt.Println("Signing the message for " + userID)
		fmt.Println("Signature Message:", signature[:50])

		if !mode.Verify(publiceKey, msg, signature) {
			panic("Signature has NOT been verified!")
		} else {
			fmt.Printf("Signature has been verified!")
		}

	} else {
		fmt.Println("Failed to get the private key for the User " + userID + " to sign the message.")
	}

}
