package aesCipher

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)


func CreateCipher(keystr string) (cipher.Block, error) {
	key, err := hex.DecodeString(keystr) //Convert string back to [] byte
	if err != nil {
		return nil, err
	}
	return aes.NewCipher(key) //create AES cipher
}
