package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)



func NewDBConnect(host string, port int, user, password, dbname, sslmode string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s  dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode))
	if err != nil {
		fmt.Println("DB ERROR")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("DB ERROR")
		return nil, err
	}
	return db, nil
}
