package repository

import (
	"auth_service"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres{
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) RegistrationPostrgres(user authservice.AuthRegistrationSerializer) (int, string, error){
	var id int
	var email string

	query := fmt.Sprintf("Insert into %s (username, first_name, last_name, email, hash_password) values ($1, $2, $3, $4, $5) returning id, email;", UserTable)

	row := r.db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Password)

	if err:= row.Scan(&id, &email); err != nil{
		return 0, "", err
	}
	return id, email, nil

	
}

func (r *AuthPostgres) ActivateUserPostgres(id int) (error) {
	query := fmt.Sprintf("Update  %s set activate = True where id=$1", UserTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r *AuthPostgres) LoginPostgres(param, value string)(authservice.LoginPostgresData, error){
	var data authservice.LoginPostgresData
	query := fmt.Sprintf("Select id, role_id, hash_password, activate, block from %s where %s = $1 limit 1;", UserTable, param)
	err := r.db.Get(&data, query, value)
	return data, err
	
}

func (r *AuthPostgres) CreateJwtRefreshPostgres(user_id, refresh string)(error){
	var id int
	query := fmt.Sprintf("Insert into %s (user_id, refresh_token) values ($1, $2) returning id;", RefreshTable)
	row := r.db.QueryRow(query, user_id, refresh)
	if err:= row.Scan(&id); err != nil{
		return err
	}
	return  nil
}

