package authservice

import (
	"time"

	 "gopkg.in/guregu/null.v3"
)


type AuthRegistrationSerializer struct{
	Username string `json:"username" binding:"required" valid:"-"`
	Email    string `json:"email" binding:"required" valid:"email"`
	Password string `json:"password" binding:"required" valid:"-"`
	FirstName string `json:"first_name" valid:"-"`
	LastName string `json:"last_name" valid:"-"`
}



type AuthRegistrationResponseSerializer  struct{
	Id int `json:"id" binding:"required" valid:"-"`
	Email string `json:"email" binding:"required" valid:"email"`

}

type ConfirmUserSerializer struct{
	Id int `json:"id" db:"id"`
}

type AuthMiddlewareSerializer struct{
	Id string `json:"id" binding:"required" valid:"-"`
	Role string `json:"role" binding:"required" valid:"-"`

}


type ProfileSerializer struct {
	Id int `json:"id" valid:"-" db:"id"`
	Username string `json:"username"  db:"username" valid:"-"`
	Email string `json:"email" valid:"email" db:"email"`
	Role string `json:"role" valid:"-" db:"role_name"`
	FirstName string `json:"first_name"  db:"first_name" valid:"-"`
	LastName string `json:"last_name"  db:"last_name" valid:"-"`
	DateTime time.Time `json:"datetime_create"  db:"datetime_create" valid:"-"`
	Image null.String `json:"image" valid:"-" db:"image"`
}


 