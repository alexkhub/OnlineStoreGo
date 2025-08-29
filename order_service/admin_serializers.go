package orderservice

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type OrderFilter struct{
	CreateAtGTE string
	CreateAtLTE null.String
	MinPrice null.Int
	MaxPrice null.Int
	PaymentMethode null.String
	Status  null.String

}

type AdminPreparatoryOrderListSerializer struct{
	Id int `db:"id"`
	UserId int64 `db:"user_id"`
	FullPrice int `db:"full_price"`
	DeliveryMethod string `db:"delivery_method"`
	PaymentStatus string `db:"payment_status"`
	CreateAt time.Time `db:"create_at"`
	Status string `db:"status"`
	DeliveryDate null.Time `db:"delivery_date"`
}

type AdminOrderListSerializer struct{
	Id int `json:"id"`
	UserId OrderUserSerializer `json:"user"`
	FullPrice int `json:"full_price"`
	DeliveryMethod string `json:"delivery_method"`
	PaymentStatus string `json:"payment_status"`
	CreateAt time.Time `json:"create_at"`
	Status string `json:"status"`
	DeliveryDate null.Time `json:"delivery_date"`
	
}



type AdminOrderStatisticSerializer struct{
	Amount null.Int `json:"amount" db:"amount"`
	TotalPrice null.Int `json:"total_price" db:"total_price"`
	AvgPrice null.Float `json:"avg_price" db:"avg_price"`

}
