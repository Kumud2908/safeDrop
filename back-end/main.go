package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	aesCipher "safeDrop/AES_cipher"
	url "safeDrop/URL"
	saveEncryptedFile "safeDrop/saveEncryptedFile"
	uniqueEncryptionKey "safeDrop/unique_Encryption_Key"
	uniqueIdentifier "safeDrop/unique_identifier"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 for local testing
	}

	// 🔥 Enable CORS
	// Disable CORS restrictions (Allow all origins)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // ✅ Allows requests from any origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"}, // ✅ Allows all headers
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 🔹 Upload route
	router.POST("/upload", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		file, err := c.FormFile("file")
		if err != nil {
			log.Println("File upload error:", err)
			c.String(http.StatusBadRequest, "Failed to retrieve file")
			return
		}

		log.Println("Uploaded file:", file.Filename)

		src, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error opening file")
			return
		}
		defer src.Close()

		// Generate unique ID and encryption key
		identifierID := uniqueIdentifier.GenerateID()
		cipherKey, err := uniqueEncryptionKey.GenerateKey()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error generating encryption key")
			return
		}

		// Create AES cipher block
		cipherBlock, err := aesCipher.CreateCipher(cipherKey)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating cipher")
			return
		}

		// Generate IV
		iv := make([]byte, cipherBlock.BlockSize())
		_, err1 := rand.Read(iv)
		if err1 != nil {
			log.Println("Error generating IV:", err1)
			return
		}

		// Encrypt and save file
		encryptedFilePath := fmt.Sprintf("%s_encrypted.dat", identifierID)

		fileURL, _, err := saveEncryptedFile.SaveEncryptedFile(cipherBlock, src, encryptedFilePath, iv)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error saving encrypted file")
			return
		}

		// Generate file download URL
		downloadURL := url.GenerateDownloadURL(identifierID, cipherKey)
		log.Println("Generated Download URL:", downloadURL)

		c.JSON(http.StatusOK, gin.H{
			"message":  fmt.Sprintf("File '%s' uploaded and encrypted as '%s'", file.Filename, encryptedFilePath),
			"file_url": fileURL,
			"id":       identifierID,
			"key":      cipherKey,
		})
		log.Println("Starting SafeDrop server on port:", port)
	})

	// 🔹 Download route (you can add this route if required)

	// Start the HTTP server
	router.Run(":" + port) // Vercel expects to listen on 8080 or a dynamically assigned port from the environment variable
}
