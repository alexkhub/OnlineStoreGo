package service

import (
	authservice "auth_service"
	"auth_service/pkg/repository"
	"bytes"
	"context"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type ProfileService struct {
	repos    repository.Profile
	minIO    *minio.Client
	producer sarama.SyncProducer
}

func NewProfileService(repos repository.Profile, minIO *minio.Client, producer sarama.SyncProducer) *ProfileService {
	return &ProfileService{
		repos:    repos,
		minIO:    minIO,
		producer: producer,
	}
}

func (s *ProfileService) UserProfile(user_id int) (authservice.ProfileSerializer, error) {

	data, err := s.repos.UserProfilePostgres(user_id)
	if err != nil {
		return data, err
	}
	image_url, err := s.minIO.PresignedGetObject(context.Background(), "user-img-bucket", data.Image.String, time.Second*24*60*60, nil)

	if err != nil {
		return data, err
	}
	data.Image.SetValid(strings.Replace(image_url.String(), "minio", "localhost", 1))
	return data, err
}

func (s *ProfileService) UpdateProfileImage(user_id int, file_data authservice.FileUploadSerializer) error {

	objectID := uuid.New().String() + file_data.FileName
	file := bytes.NewReader(file_data.Data)
	_, err := s.minIO.PutObject(context.Background(), "user-img-bucket", objectID, file, file_data.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}

	return s.repos.UpdateProfileImagePostgres(user_id, objectID)
}

func (s *ProfileService) ProfileUpdate(user_id int, user_data authservice.ProfileSerializer) error {
	return s.repos.ProfileUpdatePostgres(user_id, user_data)
}

func (s *ProfileService) ProfileDelete(user_id int) error {
	go func() {
		err := SendBlockKafkaMessageV2(s.producer, user_id)
		if err != nil {
			log.Printf("Send Block V2 Kafka %s", err.Error())
		}
	}()
	return s.repos.ProfileDeletePostgres(user_id)
}
