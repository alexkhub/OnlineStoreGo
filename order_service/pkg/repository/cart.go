package repository

import (
	"fmt"
	orderservice "order_service"

	"github.com/jmoiron/sqlx"
)

type CartPostgres struct {
	db *sqlx.DB
}

func NewCartPotgres(db *sqlx.DB) *CartPostgres{
	return &CartPostgres{
		db: db,
	}
}

func (r *CartPostgres) CartListPostgres(user_id int)([]orderservice.CartPostgresSerializer, error){
	var data []orderservice.CartPostgresSerializer

	query := fmt.Sprintf("select id, product_id, amount from %s where user_id = $1", CartTable)
	err := r.db.Select(&data, query, user_id)
	if err != nil{
		return nil, err
	}
	return data, nil

}

func (r *CartPostgres) CreateCartPostgres(user_id int, product_id int64)(orderservice.CartSerializer, error){
	var id int64

	query := fmt.Sprintf("insert into %s (product_id,  user_id) values ($1, $2) returning id", CartTable)
	row := r.db.QueryRow(query, product_id, user_id)

	if err := row.Scan(&id); err != nil {
		return orderservice.CartSerializer{}, err
	}
	return orderservice.CartSerializer{
		Id: id, 
		Amount: 1,
	}, nil
}