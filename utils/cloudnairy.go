package utils

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

func init(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")

	}
}
func UploadImage(filePath string) (string, error) {

	CLOUDINARY_URL := os.Getenv("MONGODB_URI")
	if CLOUDINARY_URL == "" {
		log.Fatal("CLOUDINARY_URL environment variable is not set")
	}

    cld, _ := cloudinary.NewFromURL(CLOUDINARY_URL)
    uploadResult, err := cld.Upload.Upload(context.Background(), filePath, uploader.UploadParams{})
    if err != nil {
        return "", err
    }
    return uploadResult.SecureURL, nil
}
