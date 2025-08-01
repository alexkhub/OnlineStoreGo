package repository

import (
	orderservice "order_service"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	CartTable = "cart"
)
type Cart interface{
	CartListPostgres(user_id int)([]orderservice.CartPostgresSerializer, error)
	CreateCartPostgres(user_id int, product_id int64)(orderservice.CartSerializer, error)
	UserCartPermissionPostgres(user_id, cart_id int ) bool
	UpdateCartPostgres(cart_id, amount int) error
	CleanCartPostgres(user_id int) error
	RemoveCartPointPostgres(cart_id int) error
}

type Order interface{

}

type Repository struct{
	Cart
	Order
}

type ReposDeps struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func NewRepository(deps ReposDeps) *Repository {
	return &Repository{
		Cart: NewCartPotgres(deps.DB),
		Order: NewOrderPostgres(deps.DB, deps.Redis),
	}
}
 