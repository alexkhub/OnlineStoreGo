package repository

import (
	// "fmt"

	authservice "auth_service"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type GRPCRepository struct {
	db *sqlx.DB
}

func NewGRPCRepository(db *sqlx.DB) *GRPCRepository {
	return &GRPCRepository{db: db}
}

func (r *GRPCRepository) GetUserDataPostgres(user_ids []int64) ([]authservice.UserDataSerializer, error) {
	var user_data []authservice.UserDataSerializer
	query := fmt.Sprintf("Select id, CONCAT(first_name, ' ', last_name) as full_name, image from %s where id = any($1) ;", UserTable)
	if err := r.db.Select(&user_data, query, pq.Array(user_ids)); err != nil {
		return nil, err
	}

	return user_data, nil
}
