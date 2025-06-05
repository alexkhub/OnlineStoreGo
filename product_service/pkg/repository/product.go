package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	productservice "product_service"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)


type ProductPostgres struct{
	db *sqlx.DB
	redisdb *redis.Client

}
func NewProductPostgres(db *sqlx.DB, redisdb *redis.Client) *ProductPostgres{
	return &ProductPostgres{
		db: db,
		redisdb: redisdb,
	}
}

func (r *ProductPostgres) CatregoListPostgres(cache bool)([]productservice.CategorySerializer, error){
	var data []productservice.CategorySerializer
	query := fmt.Sprintf("select id, name from %s;", CategoryTable)

	err := r.db.Select(&data, query)

	if err != nil{
		return data, err
	}
	go func(){
		if cache{
			cache_data, err := json.Marshal(data)
			if err != nil{
				log.Printf("cache json error %s", err )
			}
			err = r.redisdb.Set(context.Background(), RedisCategory, cache_data, time.Duration(time.Minute * 15)).Err()

			if err != nil{
				log.Printf("cache set error %s", err )
			}
		}
	}()
	return data, nil
}

func (r *ProductPostgres) ProductListPostgres()(){
	query := fmt.Sprintf(`SELECT  product.id, product.name, product.first_price, product.description,
    product.discount, product.price, product.category, image.image_uuid
	FROM product  LEFT JOIN LATERAL (
	SELECT image.image_uuid
	FROM product_image 
	JOIN image ON image.id = product_image.image
	WHERE product_image.product = product.id
	LIMIT 1) image ON true;`)
	fmt.Println(query)
}
