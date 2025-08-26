package orderservice

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type ConfirmOrderStep1Serializer struct {
	Code  int64 `json:"code"`
	Order int64
}

type OrderUserSerializer struct {
	Id       int64  `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type EmployeePreparatoryOrderDataSerializer struct {
	Id             int64     `db:"id"`
	UserId         int64     `db:"user_id"`
	FullPrice      int64     `db:"full_price"`
	PaymentMethod  string    `db:"payment_method"`
	Status         string    `db:"status"`
	DeliveryMethod string    `db:"delivery_method"`
	CreateAt       time.Time `db:"create_at"`
	DeliveryDate   null.Time `db:"delivery_date"`
	EmployeeId     null.Int  `db:"employee"`
}

type OrderPointSerializer struct {
	Id           int64  `db:"id" json:"id"`
	ProductId    int64  `db:"product_id" json:"product_id"`
	ProductPrice int64  `db:"product_price" json:"product_price"`
	ProductName  string `json:"product_name"`
	Amount       int64  `db:"amount" json:"amount"`
}

type EmployeeOrderDataSerializer struct {
	Id             int64                  `json:"id"`
	User           OrderUserSerializer    `json:"user"`
	FullPrice      int64                  `json:"full_price"`
	PaymentMethod  string                 `json:"payment_method"`
	Status         string                 `json:"status"`
	DeliveryMethod string                 `json:"delivery_method"`
	CreateAt       time.Time              `json:"create_at"`
	DeliveryDate   null.Time              `json:"delivery_date"`
	Employee       OrderUserSerializer    `json:"employee"`
	OrderPoints    []OrderPointSerializer `json:"order_points"`
}

type UpdateOrderPointSerializer struct {
	Id     int64 `json:"id" binding:"required"`
	Amount int64 `json:"amount" binding:"required,gte=0"`
}
type UpdateListOrderPointSerializer struct {
	OrderId int64
	Data    []UpdateOrderPointSerializer `json:"data" binding:"required,dive"`
}


type ConfirmOrderStep3Serializer struct {
	OrderId int64
	Employee int64
	Status string `json:"status" binding:"required"`
	PaymentStatus string `json:"payment_status"`

}