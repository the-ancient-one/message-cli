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
	"message-cli/common"
	"message-cli/config"
	"message-cli/msgcrypto"
	"os"
	"strconv"
	"time"

	"github.com/cloudflare/circl/kem/schemes"
	"github.com/cloudflare/circl/sign/dilithium"
	"github.com/spf13/cobra"
)

var message string

// sendMsgCmd represents the sendMsg command
var sendMsgCmd = &cobra.Command{
	Use:   "sendMsg",
	Short: "Send a message to a user",
	Long:  `Send message to a user. You will be prompted to enter the user ID and message to send. You need to create a user before sending a message.`,
	Run: func(cmd *cobra.Command, args []string) {

		if userID == "" {
			fmt.Print("Enter user ID: ")
			fmt.Scanln(&userID)
		}

		if !common.CheckUserExists(userID) {
			fmt.Printf("\nUser does not exist\n\n")
			slog.Error("Queried User does not exist" + userID)
			fmt.Println("Below listed are the existsing contacts.")
			ListDirectories()

			fmt.Printf("\nTo create new users, please refer the command `message-cli userID`\n\n")
			return
		} else {
			if message == "" {
				fmt.Print("Enter message: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					message = scanner.Text()
				}
			}
			slog.Info("Sending the message to the user")
			SendMsg(userID, message)
			slog.Info("Message sent to the user")
		}
	},
}

func init() {
	rootCmd.AddCommand(sendMsgCmd)
	sendMsgCmd.Flags().StringVarP(&userID, "userID", "u", "", "User ID")
	sendMsgCmd.Flags().StringVarP(&message, "message", "m", "", "Message to send")
}

func SendMsg(userID string, message string) {
	slog.Info("Hashing the message to the user" + userID)
	hashedMessage, signedMsg := hashSignMsg(userID, message)

	// Encrypt the message
	slog.Info("Encrypting the message to the user" + userID)
	sharedSecret, encryptedMessage := encryptMessage([]byte(message), userID)
	// Save the message
	slog.Info("Saving the message to the user" + userID)
	saveMessage(userID, hashedMessage, sharedSecret, signedMsg, encryptedMessage)
}

func saveMessage(userID string, hash []byte, sharedSecret []byte, signedMessage []byte, encryptedMessage []byte) {
	// Save the encrypted message
	if _, err := os.Stat("storage/" + userID + "/messages/"); os.IsNotExist(err) {
		err := os.Mkdir("storage/"+userID+"/messages/", 0755)
		if err != nil {
			fmt.Println("Failed to save message:"+userID, err)
			slog.Error("Failed to save message:" + userID + err.Error())
			return
		}
	}

	fmt.Println("Saving message...", hex.EncodeToString(encryptedMessage))
	slog.Info("Saving message to the user" + userID + hex.EncodeToString(encryptedMessage))

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
		slog.Error("Failed to marshal encrypted message to JSON:" + err.Error())
		return
	}

	msgNumber := incrementCounter(userID)

	encryptedMsgFile := "storage/" + userID + "/messages/encryptedMsg-" + strconv.Itoa(msgNumber) + ".json"
	err = os.WriteFile(encryptedMsgFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Failed to save encrypted message to file:", err)
		slog.Error("Failed to save encrypted message to file:" + err.Error())
		return
	}

	fmt.Println("Message saved successfully")
	slog.Info("Message saved successfully" + userID)
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
			slog.Error("Failed to read the Self private key file" + err.Error())
			return nil, nil
		}

		publicKeyBytes, err := os.ReadFile(pubFile)
		if err != nil {
			fmt.Println("Failed to read the Self Public key file:", err)
			slog.Error("Failed to read the Self Public key file" + err.Error())
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
			slog.Error("Signature has NOT been verified!")
			panic("Signature has NOT been verified!")
		} else {
			fmt.Println("Signature has been verified!")
			slog.Info("Signature has been verified!")
		}

		fmt.Println("Signed Message Hex len:", len(hex.EncodeToString(signedMsg)))
		slog.Info("Signed Message Hex len:" + strconv.Itoa(len(hex.EncodeToString(signedMsg))))

		return hashedMessage[:], signedMsg

	} else {
		fmt.Println("Failed to get the private key to sign the message.")
		slog.Error("Failed to get the private key to sign the message." + userID)
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
			slog.Error("Failed to read the public key file" + err.Error())
			return nil, nil
		}

		//Load the public key
		publicKey, _ := scheme.UnmarshalBinaryPublicKey([]byte(publicKeyBytes))

		ct, encryptedMessage, err := msgcrypto.Encrypt(publicKey, eseed, message)
		if err != nil {
			fmt.Println("Failed to encrypt the message:", err)
			slog.Error("Failed to encrypt the message" + err.Error())
			return nil, nil
		}

		fmt.Println("Encrypted Message Hex len:", len(hex.EncodeToString(ct)))
		return ct, encryptedMessage
	} else {
		fmt.Println("Failed to get the public key to encrypt the message.")
		slog.Error("Failed to get the public key to encrypt the message" + userID)
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
				slog.Error("Failed to create counter file" + err.Error())
				return 0
			}
			counterBytes = []byte("0")
		} else {
			fmt.Println("Failed to read counter file:", err)
			slog.Error("Failed to read counter file" + err.Error())
			return 0
		}
	}

	// Convert the counter value to an integer
	counter, err := strconv.Atoi(string(counterBytes))
	if err != nil {
		fmt.Println("Failed to convert counter value to integer:", err)
		slog.Error("Failed to convert counter value to integer" + err.Error())
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
		slog.Error("Failed to write counter file" + err.Error())
		return 0
	}

	return counter
}
