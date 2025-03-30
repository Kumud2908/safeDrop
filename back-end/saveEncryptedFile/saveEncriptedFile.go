package saveEncryptedFile

import (
	"bytes"
	"context"
	"crypto/cipher"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var (
	s3Client *s3.Client
	bucket   = "safeDropBucket"
)

func credentials() (*cloudinary.Cloudinary, context.Context, error) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, err := cloudinary.NewFromParams("djpnt1egl", "253121288254392", "J5nO4X6QcqnGO8d85gtoswJhEVU")
	if err != nil {
		return nil, nil, fmt.Errorf("Cloudinary initialization failed: %v", err)
	}

	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx, nil
}

func SaveEncryptedFile(block cipher.Block, input io.Reader, fileName string, iv []byte) (string, string, error) {
	cld, ctx, err := credentials()
	if err != nil {
		return "", "", fmt.Errorf("Cloudinary initialization error: %v", err)
	}

	// Create a buffer to hold encrypted data
	var encryptedBuffer bytes.Buffer

	// Write IV to file
	if _, err := encryptedBuffer.Write(iv); err != nil {
		fmt.Println("Error writing IV:", err)
		return "", "", err
	}

	// Create stream cipher
	stream := cipher.NewCTR(block, iv)
	if stream == nil {
		fmt.Println("Error creating cipher stream")
		return "", "", fmt.Errorf("failed to create cipher stream")
	}

	// Encrypt and write to file
	writer := &cipher.StreamWriter{S: stream, W: &encryptedBuffer}

	if _, err := io.Copy(writer, input); err != nil {
		fmt.Println("Error writing encrypted data:", err)
		return "", "", err
	}

	// Upload encrypted file to Cloudinary
	uploadResp, err := cld.Upload.Upload(ctx, bytes.NewReader(encryptedBuffer.Bytes()), uploader.UploadParams{
		PublicID:       fileName,
		UniqueFilename: false, // ✅ Fix: Use api.Bool(true) instead of plain true
		Overwrite:      true,  // ✅ Fix: Use api.Bool(true) instead of plain true
		ResourceType:   "raw", // Ensures file is stored as raw data
	})

	if err != nil {
		return "", "", fmt.Errorf("error uploading to Cloudinary: %v", err)
	}

	fmt.Println("✅ File encrypted & uploaded to Cloudinary:", uploadResp.SecureURL)
	return uploadResp.SecureURL, uploadResp.PublicID, nil
}
