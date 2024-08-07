package msgcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/schemes"
)

var meth = "Kyber512"

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
