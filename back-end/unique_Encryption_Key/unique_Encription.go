package uniqueEncriptionKey

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateKey() (string, error) {
	key := make([]byte, 32) //32bytes AES key
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil // Convert to a string for easy storage
}
