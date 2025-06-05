package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	productservice "product_service"
	"product_service/pkg/repository"
	"sync"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)


type AdminService struct{
	repos repository.Admin
	minIO *minio.Client

}

func NewAdminService(repos repository.Admin, minIO *minio.Client) *AdminService{
	return &AdminService{repos: repos, minIO: minIO}
}

func (s *AdminService) CreateCategory(data productservice.CategorySerializer) (int, error){
	return s.repos.CreateCategoryPostgres(data)
}

func (s *AdminService) CreateProduct(data productservice.AdminCreateProductSerializer) (int, error){
	return s.repos.CreateProductPostgres(data)
}

func (s *AdminService) AddImage(product int, data map[string]productservice.FileUploadSerializer) (map[string]string, error){
	file_status :=  make(map[string]string)

	if !s.repos.CheckProductIdPostgres(product){
		return file_status, errors.New("object not found")

	}
	
	urlCh := make(chan map[string]string, len(data))
	var wg sync.WaitGroup

	for _, file := range data{

		wg.Add(1)
		go func (name string, size int64, file_byte []byte){
			defer wg.Done()

			objectID := uuid.New().String() + name
			_, err := s.minIO.PutObject(context.Background(), "product",  fmt.Sprintf("product%d/%s", product, objectID), bytes.NewReader(file_byte), size, minio.PutObjectOptions{ContentType: "application/octet-stream"}) 
			if err != nil {
				urlCh <- map[string]string{
					file.FileName : "error: " + err.Error(),
				}
				return
			}
			err = s.repos.AddImagePostgres(product, objectID)
			if err != nil {
				urlCh <- map[string]string{
					file.FileName : "error: " + err.Error(),
				}
				return
			}

			urlCh <- map[string]string{
					file.FileName : "successfully",
			}

		}(file.FileName, file.Size, file.Data)	
	}
	go func() {
			wg.Wait()   
			close(urlCh) 
	}()

	for value := range urlCh{
		for key, value := range value{
			file_status[key] = value
		}
	}
	
	return file_status, nil

}
