package authservice

import (
	"gopkg.in/guregu/null.v3"
)

type AdminUserListFilter struct {
	Id        int         `json:"id" valid:"-" `
	Username  null.String `json:"username"   valid:"-"`
	Email     null.String `json:"email" valid:"email" `
	Role      null.Int    `json:"role" valid:"-"`
	FirstName null.String `json:"first_name"  valid:"-"`
	LastName  null.String `json:"last_name"   valid:"-"`
	DateTime  null.Time   `json:"datetime_create" valid:"-"`
	Actvate   null.Bool   `json:"activate" valid:"-"`
	Block     null.Bool   `json:"block" valid:"-"`
}
