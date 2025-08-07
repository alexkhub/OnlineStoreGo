package notificationsservice

import "time"


type CreateOrderKafkaMessage struct {
	Id int`json:"id"`
	User int `json:"user"`
}

type CheckOrderUUID struct {
	UserId     int       `db:"user_id"`
	OrderId  int `db:"order_id"`
	CreateTime time.Time `db:"datetime_create"`
}