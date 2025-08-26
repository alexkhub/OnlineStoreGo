package repository

import (
	"context"
	"fmt"
	productservice "product_service"

	grpc_order_service "github.com/alexkhub/OnlineStoreProto/gen/go/order_service"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)


type GRPCRepository struct {
	db *sqlx.DB
}

func NewGRPCRepository(db *sqlx.DB) *GRPCRepository {
	return &GRPCRepository{db: db}
}

func (r *GRPCRepository) GetProductCreateCartPostgres(ctx context.Context, productId int64)(*grpc_order_service.ProductDataCreateCartResponse,  error){
	var product productservice.ProductGRPCSerializer
	query := fmt.Sprintf("select id, price, name from %s where id = $1", ProductTable)

	if err := r.db.Get(&product, query, productId); err != nil{
		return nil, err
	}
	 return &grpc_order_service.ProductDataCreateCartResponse{
		Id: product.Id,
		Price: product.Price,
		Name: product.Name,
		
	}, nil

}

func (r *GRPCRepository) GetProductPostgres(ctx context.Context, productIds []int64)([]productservice.ProductGRPCSerializer, error){
	var products []productservice.ProductGRPCSerializer
	query := fmt.Sprintf("select id, price, name from %s where id = any($1);", ProductTable)
	
	if err := r.db.Select(&products, query, pq.Array(productIds)); err!= nil{
		return nil, err
	}
	return products, nil

}

func (r *GRPCRepository) GetProductPricePostgres(ctx context.Context, productIds []int64) ([]productservice.ProductPriceGRPCSerializer, error){
	var products []productservice.ProductPriceGRPCSerializer

	query := fmt.Sprintf("select id, price from %s where id = any($1) order by id;", ProductTable)
	
	if err := r.db.Select(&products, query, pq.Array(productIds)); err!= nil{
		return nil, err
	}
	return products, nil

}


func (r *GRPCRepository) GetProductNamePostgres(ctx context.Context, productIds []int64) ([]productservice.ProductNameGRPCSerializer, error){
	var products []productservice.ProductNameGRPCSerializer

	query := fmt.Sprintf("select id, name from %s where id = any($1) order by id;", ProductTable)
	
	if err := r.db.Select(&products, query, pq.Array(productIds)); err!= nil{
		return nil, err
	}
	return products, nil

}