/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/the-ancient-one/message-cli/common"
	"github.com/the-ancient-one/message-cli/config"

	"github.com/cloudflare/circl/kem/schemes"
	"github.com/cloudflare/circl/sign/dilithium"
	"github.com/spf13/cobra"
)

// For slog logger check the root.go file in the same directory

// userIDCmd represents the userID command
var userIDCmd = &cobra.Command{
	Use:   "userID",
	Short: "Create a new user/contact",
	Long:  `Create a new user/contact which will generate the pair of keys needed to communicate with the user.`,
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
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				message = scanner.Text()
			}
			slog.Info("Sending the message to the user")
			// Send the message
			SendMsg(userID, message)
			slog.Info("Completed the message sending")
		} else {
			fmt.Println("Message not sent.")
			slog.Error("Message not sent")
		}

	},
}

func init() {
	rootCmd.AddCommand(userIDCmd)

}

func checkUser(userID string) {
	if _, err := os.Stat("storage/" + userID); os.IsNotExist(err) {
		fmt.Println("User" + userID + " does not exist")
		err := os.Mkdir("storage/"+userID, 0755)
		if err != nil {
			fmt.Println("Failed to create User:"+userID, err)
			slog.Error("Failed to create User:" + userID)
			return
		}
		fmt.Println("User" + userID + " created successfully")
		slog.Info("User" + userID + " created successfully")
	} else {
		fmt.Println("User " + userID + " exists")
	}

	// Check if the user has a keys directory
	checkKeysDir(userID)

	before_mem, before_cpu, _ := common.GetSystemStats()
	// Check if the user has a signing key pair
	checkSignKeysPK(userID)
	after_mem, after_cpu, _ := common.GetSystemStats()

	used_cpu := float64(after_cpu.Total - before_cpu.Total)
	used_mem := float64(after_mem.Total - before_mem.Total)
	fmt.Printf("cpu user: %f %%\n", float64(after_cpu.User-before_cpu.User)/used_cpu*100)
	fmt.Printf("mem user: %f %%\n", float64(after_mem.Used-before_mem.Used)/used_mem*100)

	// Check if the user has a KEM key pair
	checkKEMKeysPK(userID)

	// Check if the user has a symmetric key pair for encryption
	// checkKeysSK(userID)
}

func checkKeysDir(userID string) {
	// Create the keys directory if it doesn't exist
	if _, err := os.Stat("storage/" + userID + "/keys"); os.IsNotExist(err) {
		err := os.Mkdir("storage/"+userID+"/keys", 0755)
		if err != nil {
			fmt.Println("Failed to create keys directory for User:"+userID, err)
			slog.Error("Failed to create keys directory for User:" + userID)
			return
		}
	}
}

// Check if the user has a Assymetric key pair
func checkSignKeysPK(userID string) {

	// Check if the Dilithium5 key pair already exists
	if _, err := os.Stat("storage/" + userID + "/keys/sign/privateKeySK"); !os.IsNotExist(err) {
		fmt.Println("Signing Key pair already exists for " + userID)
		slog.Info("Signing Key pair already exists for " + userID)
		return
	} else { // else create a new key pair
		modename := config.SignMode()

		// Generate Dilithium5 key pair
		mode := dilithium.ModeByName(modename)

		publicKey, privateKey, err := mode.GenerateKey(nil)
		if err != nil {
			fmt.Println("Failed to generate Key pair:", err)
			slog.Error("Failed to generate Key pair:" + err.Error())
			return
		}

		// Create the keys directory if it doesn't exist
		if _, err := os.Stat("storage/" + userID + "/keys/sign"); os.IsNotExist(err) {
			err := os.Mkdir("storage/"+userID+"/keys/sign", 0755)
			if err != nil {
				fmt.Println("Failed to create keys directory for User:"+userID, err)
				slog.Error("Failed to create keys directory for User:" + userID)
				return
			}
		}

		// Save the key pair to files
		err = os.WriteFile("storage/"+userID+"/keys/sign/privateKeySK", privateKey.Bytes(), 0644)
		if err != nil {
			fmt.Println("Failed to save private key:", err)
			slog.Error("Failed to save private key:" + err.Error())
			return
		}
		err = os.WriteFile("storage/"+userID+"/keys/sign/publicKeySK", publicKey.Bytes(), 0644)
		if err != nil {
			fmt.Println("Failed to save public key:", err)
			slog.Error("Failed to save public key:" + err.Error())
			return
		}

		fmt.Println("Signing Key pair created successfully for " + userID)
		slog.Info("Signing Key pair created successfully for " + userID)
	}
}

func checkKEMKeysPK(userID string) {
	// Check if the KEM key pair already exists
	if _, err := os.Stat("storage/" + userID + "/keys/kem/privateKeyKEM"); !os.IsNotExist(err) {
		fmt.Println("KEM Key pair already exists for " + userID)
		slog.Info("KEM Key pair already exists for " + userID)
		return
	} else { // else create a new key pair
		meth := config.KemMode()

		// Generate Kyber512 Scheme key pair
		scheme := schemes.ByName(meth)

		kseed := make([]byte, scheme.SeedSize())
		eseed := make([]byte, scheme.EncapsulationSeedSize())

		pk, sk := scheme.DeriveKeyPair(kseed)
		publicKey, _ := pk.MarshalBinary()
		privateKey, _ := sk.MarshalBinary()
		cipherText, sharedSecretSender, err := scheme.EncapsulateDeterministically(pk, eseed)
		if err != nil {
			panic(err)
		}
		sharedSecretReciever, err := scheme.Decapsulate(sk, cipherText)
		if err != nil {
			panic(err)
		}

		fmt.Println("Shared sharedSecretSender testing:", sharedSecretSender)
		slog.Info("Shared sharedSecretSender testing:" + string(sharedSecretSender))
		fmt.Println("Shared sharedSecretReciever testing:", sharedSecretReciever)
		slog.Info("Shared sharedSecretReciever testing:" + string(sharedSecretReciever))

		// Create the keys directory if it doesn't exist
		if _, err := os.Stat("storage/" + userID + "/keys/kem"); os.IsNotExist(err) {
			err := os.Mkdir("storage/"+userID+"/keys/kem", 0755)
			if err != nil {
				fmt.Println("Failed to create keys directory for User:"+userID, err)
				slog.Error("Failed to create keys directory for User:" + userID)
				return
			}
		}

		// Save the key pair to files
		err = os.WriteFile("storage/"+userID+"/keys/kem/privateKeyKEM", privateKey, 0644)
		if err != nil {
			fmt.Println("Failed to save private key:", err)
			slog.Error("Failed to save private key:" + err.Error())
			return
		}
		err = os.WriteFile("storage/"+userID+"/keys/kem/publicKeyKEM", publicKey, 0644)
		if err != nil {
			fmt.Println("Failed to save public key:", err)
			slog.Error("Failed to save public key:" + err.Error())
			return
		}

		fmt.Println("KEM Key pair created successfully for " + userID)
		slog.Info("KEM Key pair created successfully for " + userID)
	}
}
