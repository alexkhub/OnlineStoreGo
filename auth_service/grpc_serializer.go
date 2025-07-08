package authservice

import "gopkg.in/guregu/null.v3"

type UserDataSerializer struct {
	Id       int64       `db:"id"`
	FullName string      `db:"full_name"`
	Image    null.String `db:"image"`
}
