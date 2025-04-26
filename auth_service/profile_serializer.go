package authservice

import (
	"time"

	 "gopkg.in/guregu/null.v3"
)

type ProfileSerializer struct {
	Id int `json:"id" valid:"-" db:"id"`
	Username string `json:"username"  db:"username" valid:"-"`
	Email string `json:"email" valid:"email" db:"email"`
	Role string `json:"role" valid:"-" db:"role_name"`
	FirstName null.String  `json:"first_name"  db:"first_name" valid:"-"`
	LastName null.String  `json:"last_name"  db:"last_name" valid:"-"`
	DateTime time.Time `json:"datetime_create"  db:"datetime_create" valid:"-"`
	Image null.String `json:"image" valid:"-" db:"image"`
}

type FileUploadSerializer struct {
	FileName string
	Size int64
	Data []byte

}

type AdminUserListSerializer struct{
	Id int `json:"id" valid:"-" db:"id"`
	Username string `json:"username"  db:"username" valid:"-"`
	Email string `json:"email" valid:"email" db:"email"`
	Role string `json:"role" valid:"-" db:"role_name"`
	FirstName null.String  `json:"first_name"  db:"first_name" valid:"-"`
	LastName null.String  `json:"last_name"  db:"last_name" valid:"-"`
	DateTime time.Time `json:"datetime_create"  db:"datetime_create" valid:"-"`
	Actvate bool `json:"activate"  db:"activate" valid:"-"`
	Block bool `json:"block"  db:"block" valid:"-"`
}

type RoleListSerializer struct{
	Id int `json:"id" valid:"-" db:"id"`
	Name string `json:"name" valid:"-" db:"role_name"`
}
