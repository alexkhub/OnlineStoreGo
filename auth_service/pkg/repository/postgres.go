package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)


func NewDBConnect() (*sqlx.DB, error){
	db, err :=  sqlx.Open("postgres", "host=store_db port=5432 user=root password=alex0000  dbname=store_db sslmode=disable")
	if err != nil{
		fmt.Println("DB ERROR")
		return nil, err 
	}
	
	err = db.Ping()
	if err != nil{
		fmt.Println("DB ERROR")
		return nil, err 
	}
	return db, nil
}