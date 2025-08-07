package repository

import (
	notificationsservice "notifications_service"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	VerifyEmailTable = "verifyemail"
	VerifyOrderTable = "verify_order_email"
	ConfirmOrderTable = "confirm_order_email"
)

type Email interface {
	CreateVerify(uuid string, user int) error
	ChechUUID(uuid string) (notificationsservice.CheckUUIDData,  error)
}

type Order interface{
	CreateVerifyPostgres(user_id, order_id int) (uuid.UUID, error)
	CheckUUIDPostgres(uuid string) (notificationsservice.CheckOrderUUID, error)
	CodeGenerationPostgres(order_id int)(int, error)
}

type Repository struct {
	Email
	Order
}


type ReposDebs struct {
	DB *sqlx.DB
}

func NewRepository(debs ReposDebs) *Repository {
	return &Repository{
		Email: NewEmailPostgres(debs.DB),
		Order: NewOrderPostgres(debs.DB),
	}
}
