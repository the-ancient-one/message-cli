/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"message-cli/config"
	"message-cli/msgcrypto"
	"os"
	"strconv"

	"github.com/cloudflare/circl/kem/schemes"
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
	modename := config.SignMode()

	mode := dilithium.ModeByName(modename)

	// Convert the message to bytes
	msg := []byte(message)

	// Check if the Self has a private key
	if _, err := os.Stat("storage/self/keys/sign/privateKeySK"); !os.IsNotExist(err) {

		pkFile := "storage/self/keys/sign/privateKeySK"

		pubFile := "storage/self/keys/sign/publicKeySK"

		privateKeyBytes, err := os.ReadFile(pkFile)
		if err != nil {
			fmt.Println("Failed to read the Self private key file:", err)
			return
		}

		publicKeyBytes, err := os.ReadFile(pubFile)
		if err != nil {
			fmt.Println("Failed to read the Self Public key file:", err)
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

		encryptMessage(signedMsg, userID)

	} else {
		fmt.Println("Failed to get the private key to sign the message.")
	}

}

func encryptMessage(signedMsg []byte, userID string) {

	meth := config.KemMode()

	// Generate Kyber512 Scheme key pair

	scheme := schemes.ByName(meth)
	eseed := make([]byte, scheme.EncapsulationSeedSize())
	// Load the public key of the recipient
	if _, err := os.Stat("storage/" + userID + "/keys/kem/publicKeyKEM"); !os.IsNotExist(err) {

		pubFile := "storage/" + userID + "/keys/kem/publicKeyKEM"

		publicKeyBytes, err := os.ReadFile(pubFile)
		if err != nil {
			fmt.Println("Failed to read the "+userID+" Public key file:", err)
			return
		}

		//Load the public key
		publicKey, _ := scheme.UnmarshalBinaryPublicKey([]byte(publicKeyBytes))

		ct, encryptedMessage, err := msgcrypto.Encrypt(publicKey, eseed, signedMsg)
		if err != nil {
			fmt.Println("Failed to encrypt the message:", err)
			return
		}

		fmt.Printf("Ciphertext (shared key encapsulation): %x\n", ct)
		fmt.Printf("Encrypted Message: %x\n", encryptedMessage)

		// Save the encrypted message
		if _, err := os.Stat("storage/" + userID + "/messages/"); os.IsNotExist(err) {
			err := os.Mkdir("storage/"+userID+"/messages/", 0755)
			if err != nil {
				fmt.Println("Failed to save message:"+userID, err)
				return
			}
		}

		encryptedMsg := map[string]interface{}{
			"ct":               ct,
			"encryptedMessage": encryptedMessage,
		}

		jsonData, err := json.Marshal(encryptedMsg)
		if err != nil {
			fmt.Println("Failed to marshal encrypted message to JSON:", err)
			return
		}

		encryptedMsgFile := "storage/" + userID + "/messages/encryptedMsg.json"
		err = os.WriteFile(encryptedMsgFile, jsonData, 0644)
		if err != nil {
			fmt.Println("Failed to save encrypted message to file:", err)
			return
		}

		incrementCounter(userID)

		fmt.Println("Encrypted message saved to", encryptedMsgFile)
	}
}

func incrementCounter(userID string) {
	// Read the current counter value from the file
	counterFile := "storage/" + userID + "/messages/counter.txt"
	counterBytes, err := os.ReadFile(counterFile)
	if err != nil {
		fmt.Println("Failed to read counter file:", err)
		return
	}

	// Convert the counter value to an integer
	counter, err := strconv.Atoi(string(counterBytes))
	if err != nil {
		fmt.Println("Failed to convert counter value to integer:", err)
		return
	}

	// Increment the counter
	counter++

	// Convert the counter back to bytes
	counterBytes = []byte(strconv.Itoa(counter))

	// Write the updated counter value back to the file
	err = os.WriteFile(counterFile, counterBytes, 0644)
	if err != nil {
		fmt.Println("Failed to write counter file:", err)
		return
	}

	fmt.Println("Counter incremented to", counter)
}
