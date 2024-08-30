package msgcrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/the-ancient-one/message-cli/config"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/schemes"
	"github.com/cloudflare/circl/sign/dilithium"
)

var meth = config.KemMode()

var scheme = schemes.ByName(meth)

func Encrypt(pk kem.PublicKey, seed, plaintext []byte) (ciphertext, encryptedMessage []byte, err error) {
	// Generate shared secret and ciphertext deterministically
	ct, ss, err := scheme.EncapsulateDeterministically(pk, seed)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the message using the shared secret key
	block, err := aes.NewCipher(ss)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	encryptedMessage = aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ct, encryptedMessage, nil
}

func Decrypt(sk kem.PrivateKey, ct, encryptedMessage []byte) ([]byte, error) {
	// Decapsulate the shared secret key
	ss, err := scheme.Decapsulate(sk, ct)
	if err != nil {
		return nil, err
	}

	// Decrypt the message using the shared secret key
	block, err := aes.NewCipher(ss)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedMessage) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedMessage[:nonceSize], encryptedMessage[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func VerifySig(msg []byte, signedMsg []byte) (bool, error) {
	modename := config.SignMode()
	mode := dilithium.ModeByName(modename)

	pubFile := "storage/self/keys/sign/publicKeySK"
	publicKeyBytes, err := os.ReadFile(pubFile)
	if err != nil {
		fmt.Println("Failed to read the Self Public key file:", err)
		return false, err
	}

	// Load the public key
	publicKey := mode.PublicKeyFromBytes(publicKeyBytes)

	if !mode.Verify(publicKey, msg, signedMsg) {
		fmt.Println("Signature has NOT been verified!")
		return false, nil
	} else {
		return true, nil
	}
}

func VerifyHash(msg []byte, hash []byte) bool {
	hashedMessage := sha256.Sum256(msg)
	if !bytes.Equal(hashedMessage[:], hash) {
		panic("Hash has NOT been verified!")
	} else {
		return true
	}
}
