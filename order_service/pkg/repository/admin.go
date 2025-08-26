package repository

import (
	"fmt"

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


