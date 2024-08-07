/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"message-cli/config"
	"message-cli/msgcrypto"
	"os"

	"github.com/cloudflare/circl/kem/schemes"
	"github.com/spf13/cobra"
)

// readMsgCmd represents the readMsg command
var readMsgCmd = &cobra.Command{
	Use:   "readMsg",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var userID string

		fmt.Print("Enter user ID: ")
		fmt.Scanln(&userID)

		fmt.Println("Decrypting the message... ")
		decryptMessage(userID)

	},
}

func init() {
	rootCmd.AddCommand(readMsgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readMsgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readMsgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		encryptedMsgFile := "storage/" + userID + "/messages/encryptedMsg.json"
		jsonData, err := os.ReadFile(encryptedMsgFile)
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

		fmt.Println("Shared Secret Hex Len:", len(hex.EncodeToString(sharedSecret)))

		decryptedCt, err := msgcrypto.Decrypt(privateKey, []byte(sharedSecret), []byte(encryptedMessage))
		if err != nil {
			fmt.Println("Failed to decrypt the message:", err)
			return
		}

		fmt.Println("Decrypted Message:", string(decryptedCt))

		// Verify the signature
		msgcrypto.VerifySig([]byte(decryptedCt), []byte(signature))

		// Verify the hash
		msgcrypto.VerifyHash([]byte(decryptedCt), []byte(hash))
	}
}
