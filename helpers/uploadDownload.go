package helpers

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadImageToMinio(imageData *multipart.FileHeader, objectName string) error {
	// Implementasi koneksi ke MinIO
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretKey := os.Getenv("MINIO_SECRET_ACCESS_KEY_ID")
	bucketName := "gin-api"

	fileData, err := imageData.Open()
	if err != nil {
		return err
	}
	defer fileData.Close()

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Ganti menjadi true jika menggunakan koneksi aman (HTTPS)
	})
	if err != nil {
		return err
	}

	// Upload file content ke MinIO
	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, fileData, imageData.Size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
