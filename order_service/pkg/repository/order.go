package repository

import (
	"context"
	"fmt"
	"strings"
	orderservice "order_service"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)


type OrderPostgres struct{
	db      *sqlx.DB
	redisDB     *redis.Client
	gRPCProduct grpc_order_service.ProductClient
}

func NewOrderPostgres(db *sqlx.DB, redisDB *redis.Client, gRPCProduct grpc_order_service.ProductClient ) *OrderPostgres{
	return &OrderPostgres{
		db: db,
		redisDB: redisDB,
		gRPCProduct: gRPCProduct,
	}
}



func (r *OrderPostgres)PaymentMethodeListPostgres()([]orderservice.PaymentMethodeSerializer, error){
	var data []orderservice.PaymentMethodeSerializer

	query := fmt.Sprintf("select id, name, description from %s;", PaymentMethodeTable)
	if  err := r.db.Select(&data, query); err != nil{
		return nil, err
	}
	return data, nil
}


func (r *OrderPostgres) CreateOrderPostgres(order_data orderservice.CreateOrderSerializer) (int, error){
	var id int 
	var cartData []orderservice.CreateOrderCartDataSerializer
	var price int64
	args := make([]interface{}, 0)
	placeholders := []string{}



	query := fmt.Sprintf("select product_id, amount from %s  where user_id = $1 order by product_id;", CartTable)

	if err := r.db.Select(&cartData, query, order_data.User); err != nil{
		return 0, err
	}

	if len(cartData) == 0{
		return 0, fmt.Errorf("cart is empty")
	}

	productListId := make([]int64, len(cartData))

	for indx := range cartData {
		productListId[indx] = cartData[indx].ProductId
	}
	productData, err := r.gRPCProduct.GetProductPrice(context.Background(), &grpc_order_service.ProductIdRequest{Id: productListId})

	for indx := range len(cartData){
		price += cartData[indx].Amount * productData.Data[indx].Price
	}

	if err != nil{
		return 0, err
	}
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	
	query = fmt.Sprintf("insert into %s (payment_method, delivery_method, address, user_id, full_price) values ($1, $2, $3, $4, $5) returning id;", OrderTable)
	row  := tx.QueryRow(query, order_data.PaymentMethod, order_data.DeliveryMethod, order_data.Address, order_data.User, price) 

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	

	query = fmt.Sprintf("insert into  %s (product_id, product_price, amount) values", OrderPointTable)

	for i := range len(cartData) {
		n := i*3 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", n, n+1, n+2))
		args = append(args, cartData[i].ProductId, productData.Data[i].Price, cartData[i].Amount)
	}

	query += strings.Join(placeholders, ", ") + " RETURNING id"

	rows, err := tx.Query(query, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer rows.Close()

	var ids []interface{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return 0, err
		}
		ids = append(ids, id)
	}
	query = fmt.Sprintf("insert into  %s (user_order,  order_point) values", OrderOrderPointTable)
	for i := range len(ids) {
		query += fmt.Sprintf("(%d, $%d)",id, i+1)
		if i < len(ids)-1 {
			query += ", "
		}

	}
	_, err = tx.Exec(query, ids...)
	if err != nil{
		tx.Rollback()
		return 0, err
	}

	query = fmt.Sprintf("delete from %s where user_id = $1;", CartTable)

	_, err = tx.Exec(query, order_data.User)
	if err != nil{
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return id, nil
}