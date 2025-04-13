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

func (r *AuthPostgres) CreateJwtRefreshPostgres(user_id int , refresh string)(error){
	var id int
	query := fmt.Sprintf("Insert into %s (user_id, refresh_token) values ($1, $2) returning id;", RefreshTable)
	row := r.db.QueryRow(query, user_id, refresh)
	if err:= row.Scan(&id); err != nil{
		return err
	}
	return  nil
}

func (r *AuthPostgres) RefreshCheckUserPostgres(user_id int)(authservice.RefreshCheckUser, error){
	var data authservice.RefreshCheckUser
	query := fmt.Sprintf("Select role_id, activate, block from %s where id = $1 limit 1;", UserTable)
	err := r.db.Get(&data, query, user_id)
	return data, err
}

func (r *AuthPostgres) UpdateJwtRefreshPostres(user_id int , refresh, new_refresh string)(error){
	var id int
	var refresh_id int
	
	query := fmt.Sprintf("select id from %s where user_id = $1 and refresh_token = $2; ", RefreshTable)
	err := r.db.Get(&refresh_id, query, user_id, refresh)
	if err != nil{
		return err
	}

	tx, err := r.db.Begin()
	if err != nil{
		return err
	}

	query = fmt.Sprintf("delete from %s where user_id = $1 and refresh_token = $2; ", RefreshTable)
	_, err = tx.Exec(query, user_id, refresh)
	if err != nil{
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("Insert into %s (user_id, refresh_token) values ($1, $2) returning id;", RefreshTable)

	row := tx.QueryRow(query, user_id, new_refresh)
	if err:= row.Scan(&id); err != nil{
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return  err

	
}

func (r *AuthPostgres) DeleteRefreshJWTTokenPostgres(refresh string)(error){
	query := fmt.Sprintf("delete from %s where  refresh_token = $1; ", RefreshTable)
	_, err := r.db.Exec(query, refresh)
	return err
}


func (r *AuthPostgres) CloseAllSessionsPostgres(id int) (error){
	query := fmt.Sprintf("delete from %s where  user_id = $1; ", RefreshTable)
	_, err := r.db.Exec(query, id )
	return err

}
