package orderservice

import "gopkg.in/guregu/null.v3"



type PaymentMethodeSerializer struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Description null.String `json:"description" db:"description"`
}

type CreateOrderSerializer struct {
	PaymentMethod int `json:"payment_method" binding:"required" valid:"-"`
	DeliveryMethod string `json:"delivery_method" binding:"required" valid:"-"`
	Address string `json:"address" binding:"required" valid:"-"`
	User  int 
}

type CreateOrderCartDataSerializer struct {
	ProductId int64 `db:"product_id"`
	Amount int64 `db:"amount"`
}

type CreateOrderKafkaMessage struct {
	Id int`json:"id"`
	User int `json:"user"`
}