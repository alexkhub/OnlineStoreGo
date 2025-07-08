package repository

import (
	"fmt"
	productservice "product_service"

	"github.com/jmoiron/sqlx"
)

type CommentPostgres struct {
	db *sqlx.DB
}

func NewCommentPostgres(db *sqlx.DB) *CommentPostgres {
	return &CommentPostgres{db: db}
}

func (r *CommentPostgres) CreateCommentPostgres(data productservice.CreateCommentSerializer, product_id, user_id int) (int, error) {
	var id int
	query := fmt.Sprintf("insert into %s (title, message, raiting, user_id, product) values ($1, $2, $3, $4, $5) returning id;", CommentTable)
	row := r.db.QueryRow(query, data.Title, data.Message, data.Raiting, user_id, product_id)
	if err := row.Scan(&id); err != nil {

		return 0, err
	}
	return id, nil
}

func (r *CommentPostgres) RemoveUserCommentPostgres(user_id int) error {
	query := fmt.Sprintf("delete from %s where user_id = $1;", CommentTable)
	_, err := r.db.Exec(query, user_id)
	return err
}

func (r *CommentPostgres) CommentListPostgres(product_id int) ([]productservice.ListCommentPostgresSerializer, error){
	var commentList []productservice.ListCommentPostgresSerializer

	query := fmt.Sprintf("select id, title, user_id,  message, raiting, create_at from %s where product = $1 order by create_at DESC", CommentTable)

	err := r.db.Select(&commentList, query, product_id)
	if err != nil{
		return nil, err
	}
	return commentList, nil
}

func (r *CommentPostgres) RemoveCommentPostgres( comment_id int, user_id int) error {
	query := fmt.Sprintf("delete from %s where id = $1 and user_id = $2;", CommentTable)
	_, err := r.db.Exec(query, comment_id, user_id)
	return err
}