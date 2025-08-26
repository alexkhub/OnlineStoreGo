package orderservice

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type PaymentMethodeSerializer struct {
	Id          int         `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description null.String `json:"description" db:"description"`
}

type CreateOrderSerializer struct {
	PaymentMethod  int    `json:"payment_method" binding:"required" valid:"-"`
	DeliveryMethod string `json:"delivery_method" binding:"required" valid:"-"`
	Address        string `json:"address" binding:"required" valid:"-"`
	User           int
}

type CreateOrderCartDataSerializer struct {
	ProductId int64 `db:"product_id"`
	Amount    int64 `db:"amount"`
}

type CreateOrderKafkaMessage struct {
	Id   int `json:"id"`
	User int `json:"user"`
}

type UserOrderListSerializer struct{
	Id int `json:"id" db:"id"`
	FullPrice int `json:"full_price" db:"full_price"`
	DeliveryMethod string `json:"delivery_method" db:"delivery_method"`
	PaymentStatus string `json:"payment_status" db:"payment_status"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
	DeliveryDate null.Time `json:"delivery_date" db:"delivery_date"`

}


type OrderPermission struct {
	OrderId int64
	UserId int64
}

type MyError struct {
	Error error
	Code int
}