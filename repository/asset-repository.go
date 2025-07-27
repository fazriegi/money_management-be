package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
)

type IRepository interface {
	GetAssets() (result []model.Asset, err error)
}

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetAssets(req model.AssetRequest) (result []model.Asset, err error) {
	db := config.GetDatabase()

	if req.Sort == nil || *req.Sort == "" {
		order := "id ASC" // Default sorting if not provided
		req.Sort = &order
	}

	dataset := goqu.From("assets")

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
	// for row.Next() {
	// 	var asset model.Asset
	// 	if err := row.Scan(&asset); err != nil {
	// 		return nil, fmt.Errorf("failed to scan row: %w", err)
	// 	}
	// 	result = append(result, asset)
	// }

	return
}
