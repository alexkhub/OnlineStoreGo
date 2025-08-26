package repository

import (
	"fmt"
	notificationsservice "notifications_service"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) CreateVerifyPostgres(user_id, order_id int) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("insert into %s (user_id, order_id) values ($1, $2) returning id;", VerifyOrderTable)

	row := r.db.QueryRow(query, user_id, order_id)

	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil

}

func (r *OrderPostgres) CheckUUIDPostgres(uuid string) (notificationsservice.CheckOrderUUID, error) {
	var data notificationsservice.CheckOrderUUID
	query := fmt.Sprintf("select user_id, order_id, datetime_create from %s where id = $1", VerifyOrderTable)
	err := r.db.Get(&data, query, uuid)
	if err != nil {
		return notificationsservice.CheckOrderUUID{}, err
	}
	return data, nil
}

func (r *OrderPostgres) CodeGenerationPostgres(order_id int) (int, error) {
	var confirm_code int
	query := fmt.Sprintf("insert into %s (order_id) values ($1) returning confirm_code;", ConfirmOrderTable)
	row := r.db.QueryRow(query, order_id)

	if err := row.Scan(&confirm_code); err != nil {
		return 0, err
	}
	return confirm_code, nil

}
