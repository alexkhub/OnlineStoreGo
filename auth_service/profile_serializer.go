package authservice

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type ProfileSerializer struct {
	Id        int         `json:"id" valid:"-" db:"id"`
	Username  string      `json:"username"  db:"username"`
	Email     string      `json:"email" db:"email" binding:"email"`
	Role      string      `json:"role" db:"role_name"`
	FirstName null.String `json:"first_name"  db:"first_name"`
	LastName  null.String `json:"last_name"  db:"last_name"`
	DateTime  time.Time   `json:"datetime_create"  db:"datetime_create"`
	Image     null.String `json:"image" valid:"-" db:"image"`
}

type FileUploadSerializer struct {
	FileName string
	Size     int64
	Data     []byte
}

type AdminUserListSerializer struct {
	Id        int         `json:"id" db:"id"`
	Username  string      `json:"username" db:"username"`
	Email     string      `json:"email" db:"email" binding:"email"`
	Role      string      `json:"role" db:"role_name"`
	FirstName null.String `json:"first_name" db:"first_name"`
	LastName  null.String `json:"last_name" db:"last_name"`
	DateTime  time.Time   `json:"datetime_create"  db:"datetime_create"`
	Actvate   bool        `json:"activate"  db:"activate"`
	Block     bool        `json:"block"  db:"block"`
}

type RoleListSerializer struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"role_name"`
}
