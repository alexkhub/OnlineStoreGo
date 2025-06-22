package repository

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Email interface {
	CreateVerify(uuid string, user int) error
	ChechUUID(uuid string) (int, time.Time, error)
}

type Repository struct {
	Email
}

type ReposDebs struct {
	DB *sqlx.DB
}

func NewRepository(debs ReposDebs) *Repository {
	return &Repository{
		Email: NewEmailPostgres(debs.DB),
	}
}
