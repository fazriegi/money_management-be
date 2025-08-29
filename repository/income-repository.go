package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/jmoiron/sqlx"
)

type IIncomeRepository interface {
	GetList(req *model.GetIncomeRequest, userID uint, db *sqlx.DB) (result []model.Income, err error)
	BulkInsert(tx *sqlx.Tx, data *[]model.Income) error
	DeleteByPeriod(tx *sqlx.Tx, periodCode string, userID uint) error
}

type IncomeRepository struct {
}

func NewIncomeRepository() IIncomeRepository {
	return &IncomeRepository{}
}

func (r *IncomeRepository) GetList(req *model.GetIncomeRequest, userID uint, db *sqlx.DB) (result []model.Income, err error) {
	dialect := libs.GetDialect()

	if req.Sort == nil || *req.Sort == "" {
		order := "order_no" // Default sorting if not provided
		req.Sort = &order
	}

	dataset := dialect.From("incomes").Where(goqu.I("user_id").Eq(userID))

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

	result = make([]model.Income, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *IncomeRepository) BulkInsert(tx *sqlx.Tx, data *[]model.Income) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("incomes").Rows(*data)
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

func (r *IncomeRepository) DeleteByPeriod(tx *sqlx.Tx, periodCode string, userID uint) error {
	dialect := libs.GetDialect()

	dataset := dialect.Delete("incomes").Where(goqu.I("period_code").Eq(periodCode), goqu.I("user_id").Eq(userID))
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
