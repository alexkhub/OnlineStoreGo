package repository

import (
	authservice "auth_service"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type AdminPostgres struct{
	db *sqlx.DB
}

func NewAdminPostgres(db *sqlx.DB) *AdminPostgres{
	return &AdminPostgres{db: db}
}

func (r *AdminPostgres) UserListPostgres(filter url.Values)([]authservice.AdminUserListSerializer, error){
	var user_list []authservice.AdminUserListSerializer
	setValue := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if filter.Get("username") != ""{
		setValue = append(setValue, fmt.Sprintf("users.username Like $%d", argId))
		args = append(args, fmt.Sprintf("%%%s%%", filter.Get("username")))
		argId ++
	}

	if filter.Get("email") != ""{
		setValue = append(setValue, fmt.Sprintf("users.email Like $%d" , argId))
		args = append(args, fmt.Sprintf("%%%s%%", filter.Get("email")))
		argId ++
	}
	if filter.Get("first_name") != ""{
		setValue = append(setValue, fmt.Sprintf("users.first_name Like $%d" , argId))
		args = append(args, fmt.Sprintf("%%%s%%", filter.Get("first_name")))
		argId ++
	}
	if filter.Get("last_name") != ""{
		setValue = append(setValue, fmt.Sprintf("users.last_name Like $%d" , argId))
		args = append(args, fmt.Sprintf("%%%s%%", filter.Get("last_name")))
		argId ++
	}
	if filter.Get("role") != ""{
		setValue = append(setValue, fmt.Sprintf("users.role_id = $%d" , argId))
		role, err := strconv.Atoi(filter.Get("role"))
		if err != nil{
			return user_list, errors.New("role is int value")
		}
		
		args = append(args, role)
		argId ++
	}

	if filter.Get("datetime_create_gte") != ""{
		setValue = append(setValue, fmt.Sprintf("users.datetime_create > $%d" , argId))
		datetime_create_gte, err  := time.Parse(time.DateOnly, filter.Get("datetime_create_gte"))

		if err != nil{
			return user_list, errors.New("datetime create gte bad format")
		}

		
		args = append(args, datetime_create_gte)
		argId ++
	}

	if filter.Get("datetime_create_lte") != ""{
		setValue = append(setValue, fmt.Sprintf("users.datetime_create < $%d" , argId))
		datetime_create_gte, err  := time.Parse(time.DateOnly, filter.Get("datetime_create_lte"))

		if err != nil{
			return user_list, errors.New("datetime create lte bad format")
		}

		
		args = append(args, datetime_create_gte)
		argId ++
	}

	if filter.Get("activate") != ""{
		setValue = append(setValue, fmt.Sprintf("users.activate = $%d" , argId))
		activate, err  := strconv.ParseBool(filter.Get("activate"))
		if err != nil{
			return user_list, errors.New("activate is not boolean")
		}

		
		args = append(args, activate)
		argId ++
	}

	if filter.Get("block") != ""{
		setValue = append(setValue, fmt.Sprintf("users.block = $%d" , argId))
		activate, err  := strconv.ParseBool(filter.Get("block"))
		if err != nil{
			return user_list, errors.New("block is not boolean")
		}

		
		args = append(args, activate)
		argId ++
	}

	setQuery := strings.Join(setValue, " and ")
	
	
	var query string
	if argId != 1{
		query = fmt.Sprintf("select users.id, users.email,  users.username,  users.first_name, users.last_name, roles.role_name as role_name, users.datetime_create, users.activate, users.block from %s join %s on users.role_id=roles.id  where %s" , UserTable, RoleTable, setQuery)	
	}else{
		query = fmt.Sprintf("select users.id, users.email,  users.username,  users.first_name, users.last_name, roles.role_name as role_name, users.datetime_create, users.activate, users.block from %s join %s on users.role_id=roles.id" , UserTable, RoleTable)	
	}
	err := r.db.Select(&user_list, query, args... )
	if err != nil{
		return user_list, err
	}
	return user_list, nil

}

func (r *AdminPostgres) RoleListPostgres()([]authservice.RoleListSerializer, error){
	var role_list []authservice.RoleListSerializer
	query := fmt.Sprintf("select id, role_name from  %s", RoleTable )
	err := r.db.Select(&role_list, query)
	if err != nil{
		return role_list, err
	}
	return role_list, nil
}


func (r *AdminPostgres) UserBlockPostgres(user_id int)(error){
	query := fmt.Sprintf("Update %s set  block = true where id =$1", UserTable)
	_, err := r.db.Exec(query, user_id )

	return err 
}

func (r *AdminPostgres) UserUnblockPostgres(user_id int)(error){
	query := fmt.Sprintf("Update %s set  block = false where id =$1", UserTable)
	_, err := r.db.Exec(query, user_id )

	return err 
}

func (r *AdminPostgres) GetBlockDataPostgres(user_id int)( authservice.UserBlockResponseSerializer, error){
	var data authservice.UserBlockResponseSerializer
	query := fmt.Sprintf("Select email, block from %s where id = $1", UserTable)
	err := r.db.Get(&data, query, user_id) 
	return data, err 
}