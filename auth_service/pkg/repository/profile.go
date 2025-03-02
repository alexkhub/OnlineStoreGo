package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type ProfilePostgres struct{
	db *sqlx.DB
	minIO *minio.Client	
}

func NewProfilePostgres(db *sqlx.DB,  minIO *minio.Client) *ProfilePostgres{
	return &ProfilePostgres{db: db, minIO: minIO}
}

