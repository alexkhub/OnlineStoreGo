package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	productservice "product_service"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type ProductPostgres struct {
	db      *sqlx.DB
	redisdb *redis.Client
	minIO   *minio.Client
}

func NewProductPostgres(db *sqlx.DB, redisdb *redis.Client, minIO *minio.Client) *ProductPostgres {
	return &ProductPostgres{
		db:      db,
		redisdb: redisdb,
		minIO:   minIO,
	}
}

func (r *ProductPostgres) CategoryListPostgres(cache bool) ([]productservice.CategorySerializer, error) {
	var data []productservice.CategorySerializer
	query := fmt.Sprintf("select id, name from %s;", CategoryTable)

	err := r.db.Select(&data, query)

	if err != nil {
		return data, err
	}
	go func() {
		if cache {
			cache_data, err := json.Marshal(data)
			if err != nil {
				log.Printf("cache json error %s", err)
			}
			err = r.redisdb.Set(context.Background(), RedisCategory, cache_data, time.Duration(time.Minute*15)).Err()

			if err != nil {
				log.Printf("cache set error %s", err)
			}
		}
	}()
	return data, nil
}

func (r *ProductPostgres) ProductListPostgres() ([]productservice.ProductListSerailizer, error) {
	var data []productservice.ProductListSerailizer
	query := fmt.Sprintf(`SELECT  product.id, product.name, product.first_price, 
    product.discount, product.price, category.name as category, image.image_uuid
	FROM %s  left join %s on product.category = category.id LEFT JOIN LATERAL (
	SELECT image.image_uuid
	FROM %s 
	JOIN %s ON image.id = product_image.image
	WHERE product_image.product = product.id
	LIMIT 1) image ON true;`, ProductTable, CategoryTable, ProductImageTable, ImageTable)
	err := r.db.Select(&data, query)

	return data, err
}

func (r *ProductPostgres) CheckProductPostgres(id int) bool {
	var prod_id int
	query := fmt.Sprintf("Select id from %s where id = $1", ProductTable)
	err := r.db.Get(&prod_id, query, id)
	return err == nil
}

func (r *ProductPostgres) ProductDetailPostgres(id int) (productservice.ProductDetailSerailizer, error) {
	var product productservice.ProductDetailSerailizer
	query := fmt.Sprintf(`SELECT product.id, product.name, product.first_price, 
    product.discount, product.price, category.name as category 
	FROM %s left join %s on product.category = category.id where product.id = $1`, ProductTable, CategoryTable)

	err := r.db.Get(&product, query, id)
	if err != nil {
		return productservice.ProductDetailSerailizer{}, err
	}
	return product, nil
}
