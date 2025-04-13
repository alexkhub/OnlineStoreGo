package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"auth_service"
	
)

const (
	UserTable = "users"
	RoleTable = "roles"
	RefreshTable = "refresh"
)

type Authorization interface {
	RegistrationPostrgres(user authservice.AuthRegistrationSerializer) (int, string, error)
	ActivateUserPostgres(id int) (error)
	LoginPostgres(param, value string)(authservice.LoginPostgresData, error)
	CreateJwtRefreshPostgres(user_id int, refresh string)(error)
	RefreshCheckUserPostgres(user_id int)(authservice.RefreshCheckUser, error)
	UpdateJwtRefreshPostres(user_id int , refresh, new_refresh string)(error)
	DeleteRefreshJWTTokenPostgres(refresh string)(error)
	CloseAllSessionsPostgres(id int) (error)

}

type Profile interface {
	UserProfilePostgres(user_id int) (authservice.ProfileSerializer, error)
	UpdateProfileImage(user_id int, image_id string) (error)
}

type Repository struct{
	Authorization
	Profile
}

type ReposDebs struct{
	DB *sqlx.DB
	MinIO  *minio.Client
  
}

func NewRepository(debs ReposDebs) *Repository{
    return &Repository{
		Authorization: NewAuthPostgres(debs.DB),
		Profile: NewProfilePostgres(debs.DB, debs.MinIO),
		
    }
}
