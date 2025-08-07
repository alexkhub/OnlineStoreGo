package repository

import (
	"fmt"
	notificationsservice "notifications_service"
	"github.com/jmoiron/sqlx"
)

type EmailPostgres struct {
	db *sqlx.DB
}

func NewEmailPostgres(db *sqlx.DB) *EmailPostgres {
	return &EmailPostgres{db: db}
}

func (r *EmailPostgres) CreateVerify(uuid string, user int) error {
	var id int
	query := fmt.Sprintf("insert into %s (user_id, verify_uuid) values ($1, $2) returning id;", VerifyEmailTable)

	row := r.db.QueryRow(query, user, uuid)

	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil

}

func (r *EmailPostgres) ChechUUID(uuid string) (notificationsservice.CheckUUIDData, error) {
	var data notificationsservice.CheckUUIDData
	query := fmt.Sprintf("select user_id, datetime_create from %s where verify_uuid=$1", VerifyEmailTable)
	err := r.db.Get(&data, query, uuid)
	if err != nil {
		return notificationsservice.CheckUUIDData{}, err
	}
	return data, nil
}


