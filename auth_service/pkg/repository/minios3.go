package repository

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOConnect() (*minio.Client, error) {
	endpoint := "localhost:9001"
	accessKeyID := "dAg7RJqmhyLDcn7P8L38"
	secretAccessKey := "dn8P1NXBDlSmNH4O8JRSSiXSqNCgF48nYPOeseuR"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
	})
	if err != nil {
			fmt.Println("MinIO ERROR")
			return nil, err
	}
	return minioClient, nil

}