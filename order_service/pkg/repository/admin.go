package repository

import (
	"fmt"
	"log"
	orderservice "order_service"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AdminPostgres struct {
	db      *sqlx.DB
	redisDB *redis.Client
}

func NewAdminPostgres(db *sqlx.DB, redisDB *redis.Client) *AdminPostgres {
	return &AdminPostgres{
		db:      db,
		redisDB: redisDB,
	}
}

func (r *AdminPostgres) RemoveCartPointPostgres(product_id int) error {

	query := fmt.Sprintf("delete from %s where product_id = $1;", CartTable)

	_, err := r.db.Exec(query, product_id)

	return err
}


func (r *AdminPostgres) OrderListPostgres(filter orderservice.OrderFilter)([]orderservice.AdminPreparatoryOrderListSerializer, error){
	var data []orderservice.AdminPreparatoryOrderListSerializer
	setValue := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1


	setValue = append(setValue, fmt.Sprintf("create_at >= $%d", argId))
	args = append(args, filter.CreateAtGTE)
	argId += 1

	if filter.CreateAtLTE.Valid{
		setValue = append(setValue, fmt.Sprintf("create_at <= $%d", argId))
		args = append(args, filter.CreateAtLTE.String)
		argId += 1
	}
	if filter.MinPrice.Valid{
		setValue = append(setValue, fmt.Sprintf("full_price >= $%d", argId))
		args = append(args, filter.MinPrice.Int64)
		argId += 1
	}

	if filter.MaxPrice.Valid{
		setValue = append(setValue, fmt.Sprintf("full_price <= $%d", argId))
		args = append(args, filter.MaxPrice.Int64)
		argId += 1
	}

	if filter.PaymentMethode.Valid{
		setValue = append(setValue, fmt.Sprintf("payment_status = $%d", argId))
		args = append(args, filter.PaymentMethode.String)
		argId += 1
	}

	if filter.Status.Valid{
		setValue = append(setValue, fmt.Sprintf("status = $%d", argId))
		args = append(args, filter.Status.String)
		argId += 1
	}

	setQuery := strings.Join(setValue, " and ")

	query := fmt.Sprintf("select id, user_id, full_price, delivery_method, status, payment_status, create_at, delivery_date from %s where %s;", OrderTable, setQuery)
	log.Println(query)
	if err := r.db.Select(&data, query, args... ); err != nil{
		return []orderservice.AdminPreparatoryOrderListSerializer{}, err 
	}
	return data, nil
	
}


func (r *AdminPostgres) OrdersStatisticPostgres(filter orderservice.OrderFilter)([]orderservice.AdminOrderStatisticSerializer, error){
	var data []orderservice.AdminOrderStatisticSerializer

	setValue := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1


	setValue = append(setValue, fmt.Sprintf("create_at >= $%d", argId))
	args = append(args, filter.CreateAtGTE)
	argId += 1

	if filter.CreateAtLTE.Valid{
		setValue = append(setValue, fmt.Sprintf("create_at <= $%d", argId))
		args = append(args, filter.CreateAtLTE.String)
		argId += 1
	}
	if filter.MinPrice.Valid{
		setValue = append(setValue, fmt.Sprintf("full_price >= $%d", argId))
		args = append(args, filter.MinPrice.Int64)
		argId += 1
	}

	if filter.MaxPrice.Valid{
		setValue = append(setValue, fmt.Sprintf("full_price <= $%d", argId))
		args = append(args, filter.MaxPrice.Int64)
		argId += 1
	}

	if filter.PaymentMethode.Valid{
		setValue = append(setValue, fmt.Sprintf("payment_status = $%d", argId))
		args = append(args, filter.PaymentMethode.String)
		argId += 1
	}

	if filter.Status.Valid{
		setValue = append(setValue, fmt.Sprintf("status = $%d", argId))
		args = append(args, filter.Status.String)
		argId += 1
	}

	setQuery := strings.Join(setValue, " and ")

	query := fmt.Sprintf(`SELECT COUNT(*) AS amount, 
						SUM(CASE WHEN payment_status = 'Paid' and status = 'Delivered' THEN full_price END) AS total_price,
						ROUND(AVG(CASE WHEN payment_status = 'Paid' and status = 'Delivered' THEN full_price END), 2) AS avg_price  
						FROM %s 
						WHERE %s;
						`, OrderTable, setQuery)	
	if err := r.db.Select(&data, query, args...); err != nil{
		return []orderservice.AdminOrderStatisticSerializer{}, err
	}
	return data, nil
}