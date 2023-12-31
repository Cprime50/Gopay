package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func ImageUploadHelper(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create cloudinary instance
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))
	log.Printf("Connecting to Cloudinary\n")
	cld.Config.URL.Secure = true
	if err != nil {
		fmt.Println("error loading cloudinary", err)
		return "", NewInternal()
	}

	//upload file
	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: os.Getenv("CLOUDINARY_UPLOAD_FOLDER")})
	if err != nil {
		return "", err
	}
	return uploadParam.SecureURL, nil
}
