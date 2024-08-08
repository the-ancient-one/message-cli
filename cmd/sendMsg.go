/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"message-cli/config"
	"message-cli/msgcrypto"
	"os"
	"strconv"
	"time"

	"github.com/cloudflare/circl/kem/schemes"
	"github.com/cloudflare/circl/sign/dilithium"
	"github.com/spf13/cobra"
)

// sendMsgCmd represents the sendMsg command
var sendMsgCmd = &cobra.Command{
	Use:   "sendMsg",
	Short: "Send a message to a user",
	Long:  `Send message to a user. You will be prompted to enter the user ID and message to send. You need to create a user before sending a message.`,
	Run: func(cmd *cobra.Command, args []string) {
		var userID string
		var message string

		fmt.Print("Enter user ID: ")
		fmt.Scanln(&userID)

		fmt.Print("Enter message: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			message = scanner.Text()
		}

		SendMsg(userID, message)
	},
}

func init() {
	rootCmd.AddCommand(sendMsgCmd)

}

func SendMsg(userID string, message string) {
	hashedMessage, signedMsg := hashSignMsg(userID, message)

	// Encrypt the message
	sharedSecret, encryptedMessage := encryptMessage([]byte(message), userID)
	// Save the message
	saveMessage(userID, hashedMessage, sharedSecret, signedMsg, encryptedMessage)
}

func saveMessage(userID string, hash []byte, sharedSecret []byte, signedMessage []byte, encryptedMessage []byte) {
	// Save the encrypted message
	if _, err := os.Stat("storage/" + userID + "/messages/"); os.IsNotExist(err) {
		err := os.Mkdir("storage/"+userID+"/messages/", 0755)
		if err != nil {
			fmt.Println("Failed to save message:"+userID, err)
			return
		}
	}

	fmt.Println("Saving message...", hex.EncodeToString(encryptedMessage))

	encryptedMsg := map[string]interface{}{
		"hash":             hex.EncodeToString(hash),
		"sharedSecret":     hex.EncodeToString(sharedSecret),
		"signature":        hex.EncodeToString(signedMessage),
		"encryptedMessage": hex.EncodeToString(encryptedMessage),
		"timestamp":        time.Now().Unix(),
	}

	jsonData, err := json.Marshal(encryptedMsg)
	if err != nil {
		fmt.Println("Failed to marshal encrypted message to JSON:", err)
		return
	}

	msgNumber := incrementCounter(userID)

	encryptedMsgFile := "storage/" + userID + "/messages/encryptedMsg-" + strconv.Itoa(msgNumber) + ".json"
	err = os.WriteFile(encryptedMsgFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to save encrypted message to file:", err)
		return
	}

	fmt.Println("Message saved successfully")
}

func hashSignMsg(userID string, message string) ([]byte, []byte) {
	fmt.Println("Sending message to", userID)

	// Hash the message
	hashedMessage := sha256.Sum256([]byte(message))

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
			return nil, nil
		}

		publicKeyBytes, err := os.ReadFile(pubFile)
		if err != nil {
			fmt.Println("Failed to read the Self Public key file:", err)
			return nil, nil
		}

		//Load the private key
		privateKey := mode.PrivateKeyFromBytes(privateKeyBytes)

		//Load the public key
		publiceKey := mode.PublicKeyFromBytes(publicKeyBytes)

		// Sign the message
		signedMsg := mode.Sign(privateKey, msg)

		// Verify the signature
		if !mode.Verify(publiceKey, msg, signedMsg) {
			panic("Signature has NOT been verified!")
		} else {
			fmt.Println("Signature has been verified!")
		}

		fmt.Println("Signed Message Hex len:", len(hex.EncodeToString(signedMsg)))

		return hashedMessage[:], signedMsg

	} else {
		fmt.Println("Failed to get the private key to sign the message.")
	}
	return nil, nil
}

func encryptMessage(message []byte, userID string) ([]byte, []byte) {

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
			return nil, nil
		}

		//Load the public key
		publicKey, _ := scheme.UnmarshalBinaryPublicKey([]byte(publicKeyBytes))

		ct, encryptedMessage, err := msgcrypto.Encrypt(publicKey, eseed, message)
		if err != nil {
			fmt.Println("Failed to encrypt the message:", err)
			return nil, nil
		}

		fmt.Println("Encrypted Message Hex len:", len(hex.EncodeToString(ct)))
		return ct, encryptedMessage
	} else {
		fmt.Println("Failed to get the public key to encrypt the message.")
		return nil, nil
	}
}

func incrementCounter(userID string) int {
	defaultCounter := 0
	// Read the current counter value from the file
	counterFile := "storage/" + userID + "/messages/counter.txt"
	counterBytes, err := os.ReadFile(counterFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Create the file and write 0
			err = os.WriteFile(counterFile, []byte(strconv.Itoa(defaultCounter)), 0644)
			if err != nil {
				fmt.Println("Failed to create counter file:", err)
				return 0
			}
			counterBytes = []byte("0")
		} else {
			fmt.Println("Failed to read counter file:", err)
			return 0
		}
	}

	// Convert the counter value to an integer
	counter, err := strconv.Atoi(string(counterBytes))
	if err != nil {
		fmt.Println("Failed to convert counter value to integer:", err)
		return 0
	}

	// Increment the counter
	counter++

	// Convert the counter back to bytes
	counterBytes = []byte(strconv.Itoa(counter))

	// Write the updated counter value back to the file
	err = os.WriteFile(counterFile, counterBytes, 0644)
	if err != nil {
		fmt.Println("Failed to write counter file:", err)
		return 0
	}

	return counter
}
