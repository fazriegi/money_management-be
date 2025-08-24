package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/jmoiron/sqlx"
)

type IExpenseRepository interface {
	GetList(req *model.GetExpenseRequest, db *sqlx.DB) (result []model.Expense, err error)
	BulkInsert(tx *sqlx.Tx, data *[]model.Expense) error
	DeleteByPeriod(tx *sqlx.Tx, periodCode string) error
}

type ExpenseRepository struct {
}

func NewExpenseRepository() IExpenseRepository {
	return &ExpenseRepository{}
}

func (r *ExpenseRepository) GetList(req *model.GetExpenseRequest, db *sqlx.DB) (result []model.Expense, err error) {
	dialect := libs.GetDialect()

	if req.Sort == nil || *req.Sort == "" {
		order := "order_no" // Default sorting if not provided
		req.Sort = &order
	}

	dataset := dialect.From("expenses")

	if req.Search != "" {
		dataset = dataset.Where(goqu.I("name").ILike("%" + req.Search + "%"))
	}

	if req.PeriodCode != "" {
		dataset = dataset.Where(goqu.I("period_code").Eq(req.PeriodCode))
	}

	dataset = libs.PaginationRequest(dataset, req.PaginationRequest)

	sql, val, err := dataset.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	row, err := db.Queryx(sql, val...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer row.Close()

	result = make([]model.Expense, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *ExpenseRepository) BulkInsert(tx *sqlx.Tx, data *[]model.Expense) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("expenses").Rows(*data)
	sql, val, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = tx.Exec(sql, val...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *ExpenseRepository) DeleteByPeriod(tx *sqlx.Tx, periodCode string) error {
	dialect := libs.GetDialect()

	dataset := dialect.Delete("expenses").Where(goqu.Ex{"period_code": periodCode})
	sql, val, err := dataset.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = tx.Exec(sql, val...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
