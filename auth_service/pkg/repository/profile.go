package repository

import (
	authservice "auth_service"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"gopkg.in/guregu/null.v3"
)

type ProfilePostgres struct {
	db    *sqlx.DB
	minIO *minio.Client
}

func NewProfilePostgres(db *sqlx.DB, minIO *minio.Client) *ProfilePostgres {
	return &ProfilePostgres{db: db, minIO: minIO}
}

func (r *ProfilePostgres) UserProfilePostgres(user_id int) (authservice.ProfileSerializer, error) {
	var data authservice.ProfileSerializer
	query := fmt.Sprintf("select users.id, users.email,  users.username,  users.first_name, users.last_name, roles.role_name as role_name, users.datetime_create, users.image from %s join %s on users.role_id=roles.id where users.id = $1", UserTable, RoleTable)

	err := r.db.Get(&data, query, user_id)
	if err != nil {
		return data, err
	}
	return data, nil

}

func (r *ProfilePostgres) UpdateProfileImagePostgres(user_id int, image_id string) error {
	var old_image null.String

	query := fmt.Sprintf("select image from %s where id = $1 limit 1;", UserTable)

	err := r.db.Get(&old_image, query, user_id)

	if err != nil {
		return err
	}
	query = fmt.Sprintf("update  %s set image = $1 where id = $2;", UserTable)

	_, err = r.db.Exec(query, image_id, user_id)

	if err != nil {
		return err
	}

	if old_image.Valid {

		err = r.minIO.RemoveObject(context.Background(), "user-img-bucket", old_image.String, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("remove img err = %s", err.Error())
		}
		log.Println(old_image.String)
	}

	return nil

}

func (r *ProfilePostgres) ProfileUpdatePostgres(user_id int, user_data authservice.ProfileSerializer) error {
	setValue := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user_data.Username != "" {
		setValue = append(setValue, fmt.Sprintf("username=$%d", argId))
		args = append(args, user_data.Username)
		argId++

	}
	if user_data.Email != "" {
		setValue = append(setValue, fmt.Sprintf("email=$%d", argId))
		args = append(args, user_data.Email)
		argId++

	}
	if user_data.FirstName.String != "" {
		setValue = append(setValue, fmt.Sprintf("first_name=$%d", argId))
		args = append(args, user_data.FirstName.String)
		argId++

	}
	if user_data.LastName.String != "" {
		setValue = append(setValue, fmt.Sprintf("last_name=$%d", argId))
		args = append(args, user_data.LastName.String)
		argId++

	}
	setQuery := strings.Join(setValue, ", ")

	query := fmt.Sprintf("Update %s set %s WHERE  id=$%d;", UserTable, setQuery, argId)
	args = append(args, user_id)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *ProfilePostgres) ProfileDeletePostgres(user_id int) error {
	query := fmt.Sprintf("Update %s set activate = FALSE where id =$1", UserTable)
	_, err := r.db.Exec(query, user_id)

	return err

}
