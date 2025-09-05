package expense

import (
	"fmt"

	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/expense/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Insert(data *model.Expense, tx *sqlx.Tx) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Insert(data *model.Expense, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("expense").Rows(*data)
	sql, val, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = tx.Exec(sql, val...)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return nil

}
