package service

import (
	"auth_service/pkg/repository"
	grpc_product_service "github.com/alexkhub/OnlineStoreProto/gen/go/comment"
)

type GRPCService struct {
	repos       repository.GRPC
	minioClient repository.MinIO
}

func NewGRPCService(repos repository.GRPC, minioClient repository.MinIO) *GRPCService {
	return &GRPCService{repos: repos, minioClient: minioClient}
}

func (s *GRPCService) GetUserData(user_ids []int64) (*grpc_product_service.UserDataResponse, error) {
	var responce grpc_product_service.UserDataResponse

	data, err := s.repos.GetUserDataPostgres(user_ids)
	if err != nil {
		return nil, err
	}

	image_names := make([]string, 0, len(data))

	for _, user := range data {
		if user.Image.Valid {
			image_names = append(image_names, user.Image.String)
		}
	}
	urls, err := s.minioClient.GetMany("user-img-bucket", image_names)
	if err != nil {
		return nil, err
	}
	responce_user_data := make([]*grpc_product_service.UserData, 0, len(data))

	for _, user := range data {

		responce_user_data = append(responce_user_data, &grpc_product_service.UserData{
			Id:       user.Id,
			FullName: user.FullName,
			Image:    urls[user.Image.String],
		})
	}
	responce.Data = responce_user_data
	return &responce, nil

}
