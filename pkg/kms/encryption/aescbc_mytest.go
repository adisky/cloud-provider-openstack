package main

import (
	"encoding/hex"
	"fmt"
	"k8s.io/cloud-provider-openstack/pkg/kms/encryption/aescbc"
)

func main() {

	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("exampleplaintext")

	ciphertext, _ := aescbc.Encrypt(plaintext, key)

	fmt.Println("cipher: %s", string(ciphertext))

	dectext, _ := aescbc.Decrypt(ciphertext, key)

	fmt.Println("plaintext: %s", string(dectext))
}
