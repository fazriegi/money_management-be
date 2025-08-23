package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/jmoiron/sqlx"
)

type IAssetRepository interface {
	GetAssets(req *model.AssetRequest, db *sqlx.DB) ([]model.Asset, error)
	BulkInsert(tx *sqlx.Tx, data *[]model.Asset) error
	DeleteByPeriod(tx *sqlx.Tx, periodCode string) error
}

type AssetRepository struct {
}

func NewAssetRepository() IAssetRepository {
	return &AssetRepository{}
}

func (r *AssetRepository) GetAssets(req *model.AssetRequest, db *sqlx.DB) (result []model.Asset, err error) {
	dialect := libs.GetDialect()

	if req.Sort == nil || *req.Sort == "" {
		order := "order_no" // Default sorting if not provided
		req.Sort = &order
	}

	dataset := dialect.From("assets")

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

	result = make([]model.Asset, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *AssetRepository) BulkInsert(tx *sqlx.Tx, data *[]model.Asset) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("assets").Rows(*data)
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

func (r *AssetRepository) DeleteByPeriod(tx *sqlx.Tx, periodCode string) error {
	dialect := libs.GetDialect()

	dataset := dialect.Delete("assets").Where(goqu.Ex{"period_code": periodCode})
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
