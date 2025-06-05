package repository

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOConnect(host string, port int, accessKeyID, secretAccessKey string, useSSL bool ) (*minio.Client, error) {
	endpoint := fmt.Sprintf("%s:%d", host, port)
	

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