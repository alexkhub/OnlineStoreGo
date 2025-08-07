package service

import (
	"auth_service"
	"auth_service/pkg/repository"
	"net/url"

	"github.com/IBM/sarama"
	grpc_notifications_service "github.com/alexkhub/OnlineStoreProto/gen/go/notifications_service"
	grpc_product_service "github.com/alexkhub/OnlineStoreProto/gen/go/product_service"
	"github.com/minio/minio-go/v7"
)

const (
	AuthTopic    = "auth_topic"
	ConfirmTopic = "confirm_topic"
	BlockTopic   = "block_topik"
	BlockTopicV2 = "block_topik_v2"
)

type Authorization interface {
	Registration(user authservice.AuthRegistrationSerializer) (authservice.AuthRegistrationResponseSerializer, error)
	ActivateUser(id int) error
	LoginUser(user authservice.LoginUser) (authservice.JWTToken, error)
	RefreshJWTToken(refresh string) (authservice.JWTToken, error)
	DeleteRefreshJWTToken(refresh string) error
	CloseAllSessions(id int) error
}

type Profile interface {
	UserProfile(user_id int) (authservice.ProfileSerializer, error)
	UpdateProfileImage(user_id int, file_data authservice.FileUploadSerializer) error
	ProfileUpdate(user_id int, user_data authservice.ProfileSerializer) error
	ProfileDelete(user_id int) error
}

type Admin interface {
	UserList(filter url.Values) ([]authservice.AdminUserListSerializer, error)
	RoleList() ([]authservice.RoleListSerializer, error)
	UserBlock(user_id int) error
	UserUnblock(user_id int) error
}

type JWTManager interface {
	CreateJwtAccess(user_id, role_id string) (string, error)
	CreateJwtRefresh(user_id string) (string, error)
	Parse(accessToken string) (authservice.AuthMiddlewareSerializer, error)
	ParseRefreshToken(refreshToken string) (int, error)
}

type GRPC interface {
	GetUserData(user_ids []int64) (*grpc_product_service.UserDataResponse, error)
	GetUserEmail(id int64)(*grpc_notifications_service.UserEmailResponse, error)
}

type Deps struct {
	Repos      *repository.Repository
	JWTManager JWTManager
	Producer   sarama.SyncProducer
	MinIO      *minio.Client
}

type Service struct {
	Authorization
	Profile
	Admin
	GRPC
}

func NewService(deps Deps) *Service {

	return &Service{
		Authorization: NewAuthService(deps.Repos.Authorization, deps.JWTManager, deps.Producer),
		Profile:       NewProfileService(deps.Repos.Profile, deps.MinIO, deps.Producer),
		Admin:         NewAdminService(deps.Repos.Admin, deps.JWTManager, deps.Producer),
		GRPC:          NewGRPCService(deps.Repos.GRPC, deps.Repos.MinIO),
	}
}
