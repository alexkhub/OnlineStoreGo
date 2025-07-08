package repository

import (
	"auth_service"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

const (
	UserTable    = "users"
	RoleTable    = "roles"
	RefreshTable = "refresh"
)

type Authorization interface {
	RegistrationPostrgres(user authservice.AuthRegistrationSerializer) (int, string, error)
	ActivateUserPostgres(id int) error
	LoginPostgres(param, value string) (authservice.LoginPostgresData, error)
	CreateJwtRefreshPostgres(user_id int, refresh string) error
	RefreshCheckUserPostgres(user_id int) (authservice.RefreshCheckUser, error)
	UpdateJwtRefreshPostres(user_id int, refresh, new_refresh string) error
	DeleteRefreshJWTTokenPostgres(refresh string) error
	CloseAllSessionsPostgres(id int) error
}

type Profile interface {
	UserProfilePostgres(user_id int) (authservice.ProfileSerializer, error)
	UpdateProfileImagePostgres(user_id int, image_id string) error
	ProfileUpdatePostgres(user_id int, user_data authservice.ProfileSerializer) error
	ProfileDeletePostgres(user_id int) error
}

type Admin interface {
	UserListPostgres(filter url.Values) ([]authservice.AdminUserListSerializer, error)
	RoleListPostgres() ([]authservice.RoleListSerializer, error)
	UserBlockPostgres(user_id int) error
	UserUnblockPostgres(user_id int) error
	GetBlockDataPostgres(user_id int) (authservice.UserBlockResponseSerializer, error)
}

type GRPC interface {
	GetUserDataPostgres(user_ids []int64) ([]authservice.UserDataSerializer, error)
}

type MinIO interface {
	GetOne(bucketName, objectID string) (string, error)
	GetMany(bucketName string, objectIDs []string) (map[string]string, error)
	PresignedListObject(bucketName, prefix string, recursive bool) ([]string, error)
	RemoveAllObjects(bucketName, prefix string, recursive bool)
	RemoveOne(bucketName, objectID string) error
}

type Repository struct {
	Authorization
	Profile
	Admin
	GRPC
	MinIO
}

type ReposDebs struct {
	DB    *sqlx.DB
	MinIO *minio.Client
}

func NewRepository(debs ReposDebs) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(debs.DB),
		Profile:       NewProfilePostgres(debs.DB, debs.MinIO),
		Admin:         NewAdminPostgres(debs.DB),
		GRPC:          NewGRPCRepository(debs.DB),
		MinIO:         NewMinioClient(debs.MinIO),
	}
}
