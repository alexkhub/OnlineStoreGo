package repository

import (
	"context"
	"fmt"
	orderservice "order_service"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EmployeePostgres struct {
	db *sqlx.DB
}

func NewEmloyeePostgres(db *sqlx.DB) *EmployeePostgres {
	return &EmployeePostgres{
		db: db,
	}
}

func (r *EmployeePostgres) UpdateOrderPointsPosrgres(confirmData orderservice.UpdateListOrderPointSerializer) error {
	idPlaceholders := make([]string, 0, len(confirmData.Data))
	tx, err := r.db.Begin()
	parameters := make([]interface{}, 0, len(confirmData.Data)*2)
	if err != nil {
		return err
	}
	caseSQL := "CASE id "
	for i := range confirmData.Data {
		caseSQL += fmt.Sprintf("WHEN $%d::int THEN $%d::int ", i*2+1, i*2+2)
		parameters = append(parameters, confirmData.Data[i].Id, confirmData.Data[i].Amount)
		idPlaceholders = append(idPlaceholders, strconv.Itoa(int(confirmData.Data[i].Id)))
	}
	caseSQL += "END"
	query := fmt.Sprintf(`
		UPDATE %s
		SET amount = %s
		WHERE id IN (%s)`,
		OrderPointTable, caseSQL, strings.Join(idPlaceholders, ","),
	)

	_, err = tx.Exec(query, parameters...)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf("UPDATE %s SET id = $1 where id = $2;", OrderTable)
	_, err = tx.Exec(query, confirmData.OrderId, confirmData.OrderId)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}


func( r *EmployeePostgres) ConfirmOrderStep3Postgres(ctx context.Context, confirmData orderservice.ConfirmOrderStep3Serializer) error{
	var query string

	if confirmData.PaymentStatus == ""{
		query = fmt.Sprintf("update %s set employee = $1, status = $2 where id = $3", OrderTable)

		_, err := r.db.ExecContext(ctx, query, confirmData.Employee, confirmData.Status, confirmData.OrderId)
		return err
	}
	query = fmt.Sprintf("update %s set employee = $1, status = $2, payment_status = $3 where id = $4", OrderTable)

	_, err := r.db.ExecContext(ctx, query, confirmData.Employee, confirmData.Status, confirmData.PaymentStatus, confirmData.OrderId)
	return err
	
}