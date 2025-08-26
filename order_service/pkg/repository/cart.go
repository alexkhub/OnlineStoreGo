package repository

import (
	"fmt"
	orderservice "order_service"

	"github.com/jmoiron/sqlx"
)

type CartPostgres struct {
	db *sqlx.DB
}

func NewCartPotgres(db *sqlx.DB) *CartPostgres {
	return &CartPostgres{
		db: db,
	}
}

func (r *CartPostgres) CartListPostgres(user_id int) ([]orderservice.CartPostgresSerializer, error) {
	var data []orderservice.CartPostgresSerializer

	query := fmt.Sprintf("select id, product_id, amount from %s where user_id = $1", CartTable)
	err := r.db.Select(&data, query, user_id)
	if err != nil {
		return nil, err
	}
	return data, nil

}

func (r *CartPostgres) CreateCartPostgres(user_id int, product_id int64) (orderservice.CartSerializer, error) {
	var id int64

	query := fmt.Sprintf("insert into %s (product_id,  user_id) values ($1, $2) returning id;", CartTable)
	row := r.db.QueryRow(query, product_id, user_id)

	if err := row.Scan(&id); err != nil {
		return orderservice.CartSerializer{}, err
	}
	return orderservice.CartSerializer{
		Id:     id,
		Amount: 1,
	}, nil
}

func (r *CartPostgres) UserCartPermissionPostgres(user_id, cart_id int) bool {
	var id int

	query := fmt.Sprintf("Select id from %s where id = $1 and user_id = $2;", CartTable)
	err := r.db.Get(&id, query, cart_id, user_id)

	return err == nil
}

func (r *CartPostgres) UpdateCartPostgres(cart_id, amount int) error {
	query := fmt.Sprintf("Update %s set amount = $1 where id = $2;", CartTable)
	_, err := r.db.Exec(query, cart_id, amount)
	return err
}

func (r *CartPostgres) CleanCartPostgres(user_id int) error {
	query := fmt.Sprintf("delete from %s where user_id = $1;", CartTable)

	_, err := r.db.Exec(query, user_id)
	return err

}

func (r *CartPostgres) RemoveCartPointPostgres(cart_id int) error {
	query := fmt.Sprintf("delete from %s where id = $1;", CartTable)

	_, err := r.db.Exec(query, cart_id)
	return err

}
