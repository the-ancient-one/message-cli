/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/sha256"
	"fmt"
	"message-cli/config"
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

	// Hash the message
	hashedMessage := sha256.Sum256([]byte(message))
	fmt.Printf("Hashed: %x\n", hashedMessage)

	// Determine the mode
	modename := "Dilithium5"

	mode := dilithium.ModeByName(modename)

	// Convert the message to bytes
	msg := []byte(message)

	// Check if the user has a private key
	if _, err := os.Stat("storage/" + userID + "/keys/sign/privateKey"); !os.IsNotExist(err) {

		pkFile := "storage/" + userID + "/keys/sign/privateKey"

		pubFile := "storage/" + userID + "/keys/sign/publicKey"

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

		//Load the private key
		privateKey := mode.PrivateKeyFromBytes(privateKeyBytes)

		//Load the public key
		publiceKey := mode.PublicKeyFromBytes(publicKeyBytes)

		// append the message to the hash
		hashedMsg := append(msg, hashedMessage[:]...)

		fmt.Println("Hash + Message:", hashedMsg)

		// Sign the message
		signedMsg := mode.Sign(privateKey, hashedMsg)

		fmt.Println("Signature Message:", signedMsg[:50])

		// Verify the signature
		if !mode.Verify(publiceKey, hashedMsg, signedMsg) {
			panic("Signature has NOT been verified!")
		} else {
			fmt.Println("Signature has been verified!")
		}

		encryptMessage(signedMsg)

	} else {
		fmt.Println("Failed to get the private key for the User " + userID + " to sign the message.")
	}

}

func encryptMessage(encryptMsg []byte) {
	// Encrypt the message
	fmt.Println("Encrypted the message: ", encryptMsg)
	fmt.Println("Symmetric password is " + config.AesPasswd())
}
