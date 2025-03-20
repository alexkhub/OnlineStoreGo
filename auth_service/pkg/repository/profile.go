package repository

import (
	authservice "auth_service"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type ProfilePostgres struct{
	db *sqlx.DB
	minIO *minio.Client	
}

func NewProfilePostgres(db *sqlx.DB,  minIO *minio.Client) *ProfilePostgres{
	return &ProfilePostgres{db: db, minIO: minIO}
}


func (r *ProfilePostgres) UserProfilePostgres (user_id int) (authservice.ProfileSerializer, error){
	var data authservice.ProfileSerializer
	query := fmt.Sprintf("select users.id,  users.username,  users.first_name, users.last_name, roles.role_name as role_name, users.datetime_create, users.image from %s join %s on users.role_id=roles.id where users.id = $1", UserTable, RoleTable)

	err := r.db.Get(&data, query, user_id)
	if err != nil{
		return data, err
	}
	return data, nil 

}