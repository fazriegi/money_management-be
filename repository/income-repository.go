package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/jmoiron/sqlx"
)

type IIncomeRepository interface {
	GetList(req *model.GetIncomeRequest) (result []model.Income, err error)
	BulkInsert(tx *sqlx.Tx, data *[]model.Income) error
	DeleteByPeriod(tx *sqlx.Tx, periodCode string) error
}

type IncomeRepository struct {
}

func NewIncomeRepository() *IncomeRepository {
	return &IncomeRepository{}
}

func (r *IncomeRepository) GetList(req *model.GetIncomeRequest) (result []model.Income, err error) {
	db := config.GetDatabase()

	if req.Sort == nil || *req.Sort == "" {
		order := "order_no" // Default sorting if not provided
		req.Sort = &order
	}

	dataset := goqu.From("incomes")

	if req.Search != "" {
		dataset = dataset.Where(goqu.Ex{"name": req.Search})
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

	result = make([]model.Income, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *IncomeRepository) BulkInsert(tx *sqlx.Tx, data *[]model.Income) error {
	dataset := goqu.Insert("incomes").Rows(*data)
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

func (r *IncomeRepository) DeleteByPeriod(tx *sqlx.Tx, periodCode string) error {
	dataset := goqu.Delete("incomes").Where(goqu.Ex{"period_code": periodCode})
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
