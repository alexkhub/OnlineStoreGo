package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)


type OrderPostgres struct{
	db      *sqlx.DB
	redisDB     *redis.Client
}

func NewOrderPostgres(db *sqlx.DB, redisDB *redis.Client ) *OrderPostgres{
	return &OrderPostgres{
		db: db,
		redisDB: redisDB,
	}
}