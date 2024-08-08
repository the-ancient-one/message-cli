/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	"github.com/spf13/cobra"
)

// readMsgCmd represents the readMsg command
var readMsgCmd = &cobra.Command{
	Use:   "readMsg",
	Short: "Access the conversation with the user",
	Long:  `Access the conversation histroy with the user and listed in table format.`,
	Run: func(cmd *cobra.Command, args []string) {
		var userID string

		fmt.Print("Enter user ID: ")
		fmt.Scanln(&userID)

		fmt.Printf("\nDecrypting the message... \n")
		decryptMessage(userID)

	},
}

func init() {
	rootCmd.AddCommand(readMsgCmd)

}

func decryptMessage(userID string) {

	meth := config.KemMode()

	scheme := schemes.ByName(meth)
	// eseed := make([]byte, scheme.EncapsulationSeedSize())
	// Load the private key of the recipient
	if _, err := os.Stat("storage/" + userID + "/keys/kem/privateKeyKEM"); !os.IsNotExist(err) {

		pubFile := "storage/" + userID + "/keys/kem/privateKeyKEM"

		privateKeyBytes, err := os.ReadFile(pubFile)
		if err != nil {
			fmt.Println("Failed to read the "+userID+" Private key file:", err)
			return
		}

		//Load the private key
		privateKey, _ := scheme.UnmarshalBinaryPrivateKey([]byte(privateKeyBytes))

		// Load the encrypted message

		fmt.Println("Reading the encrypted message files...")

		fmt.Println("Conversation:")
		fmt.Println("--------------------------------------------------")
		fmt.Printf("| %-15s |%-29s| %-17s |%-22s |%-17s  |\n", "User ID", "Timestamp", "Hash Verification", "Signature Verification", "Decrypted Message")

		encryptedMsgFiles, err := common.ListEncryptedMsgFiles(userID)
		if err != nil {
			fmt.Println("Failed to list the encrypted message files:", err)
			return
		}
		for _, file := range encryptedMsgFiles {
			jsonData, err := os.ReadFile(file)
			if err != nil {
				fmt.Println("Failed to read the encrypted message file:", err)
				return
			}

			var data map[string]interface{}
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				fmt.Println("Failed to unmarshal the encrypted message JSON:", err)
				return
			}

			sharedSecret, _ := hex.DecodeString(data["sharedSecret"].(string))
			encryptedMessage, _ := hex.DecodeString(data["encryptedMessage"].(string))
			signature, _ := hex.DecodeString(data["signature"].(string))
			hash, _ := hex.DecodeString(data["hash"].(string))
			timestamp, ok := data["timestamp"].(float64)
			if !ok || timestamp == 0 {
				timestamp = 946684800
			}
			// fmt.Println("printign time ", time.Unix(timestamp, 0))

			// fmt.Println("Shared Secret Hex Len:", len(hex.EncodeToString(sharedSecret)))

			decryptedCt, err := msgcrypto.Decrypt(privateKey, []byte(sharedSecret), []byte(encryptedMessage))
			if err != nil {
				fmt.Println("Failed to decrypt the message:", err)
				return
			}
			// Verify the signature
			verifiedHash, err := msgcrypto.VerifySig([]byte(decryptedCt), []byte(signature))
			if err != nil {
				fmt.Println("Failed to verify the signature:", err)
			}

			// Verify the hash
			verifiedSign := msgcrypto.VerifyHash([]byte(decryptedCt), []byte(hash))

			fmt.Printf("| %-15s |%-29s| %-17s |%-22s |%-17s  |\n", userID, time.Unix(int64(timestamp), 0), strconv.FormatBool(verifiedHash), strconv.FormatBool(verifiedSign), string(decryptedCt))

		}
		fmt.Println("--------------------------------------------------")
	}
}
