package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryClient struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryClient() (*CloudinaryClient, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return &CloudinaryClient{cld: cld}, nil
}

func (c *CloudinaryClient) UploadImage(file multipart.File, reportID string, longitude, latitude float64, description string) (string, error) {
	ctx := context.Background()

	contextData := api.CldAPIMap{
		"Id":          reportID,
		"description": description,
		"latitude":    fmt.Sprintf("%.6f", latitude),
		"longitude":   fmt.Sprintf("%.6f", longitude),
	}

	uploadParams := uploader.UploadParams{
		PublicID: reportID,
		Folder:   "reports",
		Context:  contextData,
	}

	result, err := c.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
