package service

import (
	"bytes"

	"context"
	"errors"
	"fmt"
	"log"

	productservice "product_service"
	"product_service/pkg/repository"
	"sync"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type AdminService struct {
	repos       repository.Admin
	minIO       *minio.Client
	minioClient repository.MinIO
	redisDB     *redis.Client
	producer   	sarama.SyncProducer
}

func NewAdminService(repos repository.Admin, minIO *minio.Client, minioClient repository.MinIO, redisDB *redis.Client, producer sarama.SyncProducer) *AdminService {
	return &AdminService{repos: repos, minIO: minIO, minioClient: minioClient, redisDB: redisDB, producer: producer}
}

func (s *AdminService) CreateCategory(data productservice.CategorySerializer) (int, error) {
	return s.repos.CreateCategoryPostgres(data)
}

func (s *AdminService) CreateProduct(data productservice.AdminCreateProductSerializer) (int, error) {
	go s.redisDB.Del(context.Background(), "products")
	return s.repos.CreateProductPostgres(data)
}

func (s *AdminService) AddImage(product int, data map[string]productservice.FileUploadSerializer) (map[string]string, error) {
	file_status := make(map[string]string)

	if !s.repos.CheckProductIdPostgres(product) {
		return file_status, errors.New("object not found")

	}
	urlCh := make(chan map[string]string, len(data))
	var wg sync.WaitGroup

	for _, file := range data {

		wg.Add(1)
		go func(name string, size int64, file_byte []byte) {
			defer wg.Done()

			objectID := uuid.New().String() + name
			_, err := s.minIO.PutObject(context.Background(), "product", fmt.Sprintf("product%d/%s", product, objectID), bytes.NewReader(file_byte), size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
			if err != nil {
				urlCh <- map[string]string{
					file.FileName: "error: " + err.Error(),
				}
				return
			}
			err = s.repos.AddImagePostgres(product, objectID)
			if err != nil {
				urlCh <- map[string]string{
					file.FileName: "error: " + err.Error(),
				}
				return
			}

			urlCh <- map[string]string{
				file.FileName: "successfully",
			}

		}(file.FileName, file.Size, file.Data)
	}
	go func() {
		wg.Wait()
		close(urlCh)
	}()

	for values := range urlCh {
		for key, value := range values {
			file_status[key] = value
		}
	}
	go s.redisDB.Del(context.Background(), fmt.Sprintf("product%d", product), "products")
	return file_status, nil

}

func (s *AdminService) ProductDelete( product_id int) error {
	go s.minioClient.RemoveAllObjects("product", fmt.Sprintf("product%d", product_id), true)
	go s.redisDB.Del(context.Background(), fmt.Sprintf("product%d", product_id), "products")

	go func(producer sarama.SyncProducer, product_id int){
		err := SendDeleteProductKafkaMessage(producer, product_id)
		if err != nil{
			log.Printf("send delete product error: %s", err.Error())
		}
	}(s.producer, product_id)

	return s.repos.DeleteProductPostgres(product_id)
}

func (s *AdminService) AdminProductDetail(id int) (productservice.AdminProductDetailSerailizer, error) {
	data, err := s.repos.AdminProductDetailPostgres(id)
	if err != nil {
		return productservice.AdminProductDetailSerailizer{}, err
	}
	images, err := s.repos.GetImagesPostgres(id)
	if err != nil {
		return productservice.AdminProductDetailSerailizer{}, err
	}

	image_names := make([]string, 0, len(images))

	for _, image := range images {
		image_names = append(image_names, fmt.Sprintf("product%d/%s", id, image.Name.String))
	}

	urls, err := s.minioClient.GetMany("product", image_names)
	if err != nil {
		return productservice.AdminProductDetailSerailizer{}, err
	}
	for indx, image := range images {
		if image.Name.Valid {
			url, ok := urls[fmt.Sprintf("product%d/%s", id, image.Name.String)]
			if ok {
				images[indx].Link.SetValid(url)
			}
		}
	}
	data.Images = images

	return data, nil

}

func (s *AdminService) RemoveImage(product_id int, name string) error {
	err := s.repos.DeleteImagePostgres(name)
	if err != nil {
		return err
	}
	go s.redisDB.Del(context.Background(), fmt.Sprintf("product%d", product_id))
	return s.minioClient.RemoveOne("product", fmt.Sprintf("product%d/%s", product_id, name))

}

func (s *AdminService) UpdateProduct(product_id int, product_data productservice.AdminUpdateProductSerializer) error {
	go s.redisDB.Del(context.Background(), fmt.Sprintf("product%d", product_id), "products")
	return s.repos.UpdateProductPostgres(product_id, product_data)
}

func (s *AdminService) RemoveComment(comment_id int) error{
	return s.repos.RemoveCommentPostgres(comment_id)
}
