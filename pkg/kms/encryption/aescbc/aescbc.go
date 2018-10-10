package aescbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Encrypt plain text
func Encrypt(data, key []byte) (ciphertext []byte, err error) {

	// NewCipher returns a new cipher block, the key argument should be AES key
	// either 16, 24 or 32 bytes to select AES-128, AES-192, AES-256, 32 byte is preferred
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// determine the padding size for block cipher
	paddingSize := aes.BlockSize - (len(data) % aes.BlockSize)
	plaintext := make([]byte, len(data)+paddingSize)
	// copy data and padding
	copy(plaintext, data)
	copy(plaintext[len(data):], bytes.Repeat([]byte{byte(paddingSize)}, paddingSize))

	// create slice to hold ciphertext, iv
	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return
}

// Decrypt plaintext
func Decrypt(data, key []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	paddingLength := int(ciphertext[len(ciphertext)-1])
	dataLength := len(ciphertext) - paddingLength
	plaintext = ciphertext[:dataLength]

	return
}
