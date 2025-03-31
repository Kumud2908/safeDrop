package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"

	aesCipher "safeDrop/AES_cipher"
	url "safeDrop/URL"
	saveEncryptedFile "safeDrop/saveEncryptedFile"
	uniqueEncryptionKey "safeDrop/unique_Encryption_Key"
	uniqueIdentifier "safeDrop/unique_identifier"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: No .env file found, using default values")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Fiber app
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))
	

	// Upload route
	app.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			log.Println("File upload error:", err)
			return c.Status(fiber.StatusBadRequest).SendString("Failed to retrieve file")
		}

		log.Println("Uploaded file:", file.Filename)

		src, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
		}
		defer src.Close()

		// Generate unique ID and encryption key
		identifierID := uniqueIdentifier.GenerateID()
		cipherKey, err := uniqueEncryptionKey.GenerateKey()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error generating encryption key")
		}

		// Create AES cipher block
		cipherBlock, err := aesCipher.CreateCipher(cipherKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating cipher")
		}

		// Generate IV
		iv := make([]byte, cipherBlock.BlockSize())
		_, err = rand.Read(iv)
		if err != nil {
			log.Println("Error generating IV:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error generating IV")
		}

		// Encrypt and save file
		encryptedFilePath := fmt.Sprintf("%s_encrypted.dat", identifierID)
		fileURL, _, err := saveEncryptedFile.SaveEncryptedFile(cipherBlock, src, encryptedFilePath, iv)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving encrypted file")
		}

		// Generate file download URL
		downloadURL := url.GenerateDownloadURL(identifierID, cipherKey)
		log.Println("Generated Download URL:", downloadURL)

		return c.JSON(fiber.Map{
			"message":  fmt.Sprintf("File '%s' uploaded and encrypted as '%s'", file.Filename, encryptedFilePath),
			"file_url": fileURL,
			"id":       identifierID,
			"key":      cipherKey,
		})
	})

	// Start Fiber server
	log.Printf("🚀 SafeDrop server starting on PORT: %s", port)
	log.Fatal(app.Listen(":" + port))
}
