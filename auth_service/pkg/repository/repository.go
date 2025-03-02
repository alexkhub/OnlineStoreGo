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
	CreateJwtRefreshPostgres(user_id, refresh string)(error)
}

type Profile interface {

}

type Repository struct{
	Authorization
	Profile
}

type ReposDebs struct{
	DB *sqlx.DB
    MinIO *minio.Client
}

func NewRepository(debs ReposDebs) *Repository{
    return &Repository{
		Authorization: NewAuthPostgres(debs.DB),
		Profile: NewProfilePostgres(debs.DB, debs.MinIO),
		
    }
}
