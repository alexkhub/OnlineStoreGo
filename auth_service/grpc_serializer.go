package authservice

import "gopkg.in/guregu/null.v3"

type UserDataSerializer struct {
	Id       int64       `db:"id"`
	FullName string      `db:"full_name"`
	Image    null.String `db:"image"`
}

type OrderUserDataSerializer struct {
	Id       int64  `db:"id"`
	FullName string `db:"full_name"`
	Email    string `db:"email"`
}
