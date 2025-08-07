package notificationsservice

import (
	"time"
)

type AuthRegistrationResponseSerializer struct {
	Id    int    `json:"id" binding:"required" valid:"-"`
	Email string `json:"email" binding:"required" valid:"email"`
}

type CheckUUIDData struct {
	UserId     int       `db:"user_id"`
	CreateTime time.Time `db:"datetime_create"`
}
type UserBlockResponseSerializer struct {
	Email string `json:"email" binding:"required" valid:"email" db:"email"`
	Block bool   `json:"block" binding:"required" valid:"email" db:"block"`
}
