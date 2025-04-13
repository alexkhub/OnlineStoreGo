package repository

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOConnect() (*minio.Client, error) {
	endpoint := "minio:9000"
	accessKeyID := "SX1o3y8kcjMpN6lLjsHh"
	secretAccessKey := "PuRrMhFaap1f3RBDd0S5pHNly9dpFYtzO2BKNg7M"
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