package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// Initialization vector for the AES algorithm
var initVector = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// EncryptString encrypts the string with given key
func EncryptString(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, initVector)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText)
}

// DecryptString decrypts the encrypted string to original
func DecryptString(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	ciphertext, _ := base64.StdEncoding.DecodeString(text)
	cfb := cipher.NewCFBEncrypter(block, initVector)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext)
}
