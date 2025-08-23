package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/jmoiron/sqlx"
)

type ILiabilityRepository interface {
	GetList(req *model.GetLiabilityRequest, db *sqlx.DB) (result []model.Liability, err error)
	BulkInsert(tx *sqlx.Tx, data *[]model.Liability) error
	DeleteByPeriod(tx *sqlx.Tx, periodCode string) error
}

type LiabilityRepository struct {
}

func NewLiabilityRepository() ILiabilityRepository {
	return &LiabilityRepository{}
}

func (r *LiabilityRepository) GetList(req *model.GetLiabilityRequest, db *sqlx.DB) (result []model.Liability, err error) {
	dialect := libs.GetDialect()

	if req.Sort == nil || *req.Sort == "" {
		order := "order_no" // Default sorting if not provided
		req.Sort = &order
	}

	dataset := dialect.From("liabilities")

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

	result = make([]model.Liability, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *LiabilityRepository) BulkInsert(tx *sqlx.Tx, data *[]model.Liability) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("liabilities").Rows(*data)
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

func (r *LiabilityRepository) DeleteByPeriod(tx *sqlx.Tx, periodCode string) error {
	dialect := libs.GetDialect()

	dataset := dialect.Delete("liabilities").Where(goqu.Ex{"period_code": periodCode})
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
