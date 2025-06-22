package repository

import (
	"context"
	"fmt"

	"log"
	productservice "product_service"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"

)

type AdminPostgres struct {
	db      *sqlx.DB
	redisdb *redis.Client
	minIO   *minio.Client
}

func NewAdminPostgres(db *sqlx.DB, redisdb *redis.Client, minIO *minio.Client) *AdminPostgres {
	return &AdminPostgres{db: db, redisdb: redisdb, minIO: minIO}
}

func (r *AdminPostgres) CreateCategoryPostgres(data productservice.CategorySerializer) (int, error) {
	var id int
	query := fmt.Sprintf("Insert into %s (name) values ($1) returning id;", CategoryTable)

	row := r.db.QueryRow(query, data.Name)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	go func() {
		err := r.redisdb.Del(context.Background(), RedisCategory).Err()
		if err != nil {
			log.Printf("cache category del error %s", err.Error())
		}
	}()
	return id, nil
}

func (r *AdminPostgres) CreateProductPostgres(data productservice.AdminCreateProductSerializer) (int, error) {
	var id int

	tx, err := r.db.Begin()
	if err != nil {

		return 0, err
	}
	query := fmt.Sprintf("insert into %s (name, first_price, discount, description, category) values($1, $2, $3, $4, $5) returning id;", ProductTable)
	row := tx.QueryRow(query, data.Name, data.Price, data.Discount, data.Description, data.Category)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, nil

}

func (r *AdminPostgres) CheckProductIdPostgres(product_id int) bool {
	var id int

	query := fmt.Sprintf("select count(id) from  %s where id = $1", ProductTable)

	err := r.db.Get(&id, query, product_id)

	if err != nil || id == 0 {
		return false
	}
	return true
}

func (r *AdminPostgres) AddImagePostgres(product_id int, image string) error {
	var image_id int
	var product_image_id int

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("insert into %s (image_uuid) values ($1) returning id;", ImageTable)
	row := tx.QueryRow(query, image)
	if err := row.Scan(&image_id); err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("insert into %s (product, image) values ($1, $2)returning id;", ProductImageTable)
	row = tx.QueryRow(query, product_id, image_id)
	if err := row.Scan(&product_image_id); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *AdminPostgres) DeleteProductPostgres(id int) error {
	query := fmt.Sprintf("delete from %s where id = $1;", ProductTable)
	_, err := r.db.Exec(query, id)
	return err
}


func (r *AdminPostgres) AdminProductDetailPostgres(id int) (productservice.AdminProductDetailSerailizer, error){
	var product productservice.AdminProductDetailSerailizer
	query := fmt.Sprintf(`SELECT product.id, product.name, product.first_price, 
    product.discount, product.price, category.name as category 
	FROM %s join %s on product.category = category.id where product.id = $1`, ProductTable, CategoryTable)

	err := r.db.Get(&product, query, id)

	 
	
	return product, err
}

func (r *AdminPostgres) GetImage(product_id int) ([]productservice.ImageSerializer, error){
	var images []productservice.ImageSerializer
	query := fmt.Sprintf(`select image.image_uuid from %s 
						join %s on product_image.image=image.id 
						where product_image.product = $1;`, ProductImageTable, ImageTable)	
	err := r.db.Select(&images, query, product_id)

	return images, err
}

func (r *AdminPostgres) DeleteImagePostgres(name string) error {
	query := fmt.Sprintf("delete from %s where image_uuid = $1;", ImageTable)
	_, err := r.db.Exec(query, name)
	return err

}