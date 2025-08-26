package repository

import (
	"context"
	"fmt"
	"log"

	notificationsservice "notifications_service"
	"time"

	"github.com/jmoiron/sqlx"
)

type GRPCPostgres struct {
	db *sqlx.DB
}

func NewGRPCPostgres(db *sqlx.DB) *GRPCPostgres {
	return &GRPCPostgres{
		db: db,
	}
}

func (r *GRPCPostgres) CheckCodePostgres(ctx context.Context, orderData notificationsservice.CheckCodeSeralizer) (time.Time, error) {
	var codeTime time.Time

	query := fmt.Sprintf("select datetime_create from %s where order_id = $1 and confirm_code = $2", ConfirmOrderTable)
	err := r.db.Get(&codeTime, query, orderData.OrderId, orderData.Code)
	if err != nil {
		return time.Time{}, err
	}
	return codeTime, nil
}

func (r *GRPCPostgres) GetUserIdPostgres(ctx context.Context, orderId int64) (int64, error) {
	var user_id int64
	query := fmt.Sprintf("select user_id from %s where order_id = $1 limit 1;", VerifyOrderTable)

	err := r.db.Get(&user_id, query, orderId)
	if err != nil {
		return 0, err
	}
	return user_id, nil
}

func (r *GRPCPostgres) GenerateNewCodePostgres(ctx context.Context, orderId int64) (int, error) {
	var confirmCode int
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	log.Println(1)
	query := fmt.Sprintf("delete from %s where order_id = $1;", ConfirmOrderTable)

	_, err = tx.Exec(query, orderId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	log.Println(2)
	query = fmt.Sprintf("insert into %s (order_id) values ($1) returning confirm_code;", ConfirmOrderTable)
	row := tx.QueryRow(query, orderId)

	if err := row.Scan(&confirmCode); err != nil {
		tx.Rollback()
		return 0, err
	}
	log.Println(3)
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return confirmCode, nil
}
