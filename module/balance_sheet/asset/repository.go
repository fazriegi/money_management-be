package asset

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/balance_sheet/asset/model"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Repository interface {
	Insert(data *model.Asset, tx *sqlx.Tx) error
	ListCategory(userID uint, db *sqlx.DB) (result []model.AssetCategory, err error)
	List(req *model.ListRequest, db *sqlx.DB) (result []model.GetAsset, total uint, err error)
	Update(userId, id uint, data map[string]any, tx *sqlx.Tx) error
	Delete(userId, id uint, tx *sqlx.Tx) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Insert(data *model.Asset, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("asset").Rows(*data)
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

func (r *repository) ListCategory(userID uint, db *sqlx.DB) (result []model.AssetCategory, err error) {
	dialect := libs.GetDialect()

	dataset := dialect.From("asset_category").Where(goqu.I("user_id").Eq(userID))

	sql, val, err := dataset.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	row, err := db.Queryx(sql, val...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer row.Close()

	result = make([]model.AssetCategory, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *repository) List(req *model.ListRequest, db *sqlx.DB) (result []model.GetAsset, total uint, err error) {
	if req.Sort == nil {
		sort := "created_at asc"
		req.Sort = &sort
	}

	dialect := libs.GetDialect()

	dataset := dialect.
		From("asset").
		Join(goqu.T("asset_category").As("ac"), goqu.On(
			goqu.I("ac.id").Eq(goqu.I("asset.category_id")),
			goqu.I("ac.user_id").Eq(goqu.I("asset.user_id")),
		)).
		Select(
			goqu.I("asset.id"),
			goqu.I("asset.category_id"),
			goqu.I("ac.name").As("category"),
			goqu.I("asset.amount"),
			goqu.I("asset.value"),
			goqu.I("asset.user_id"),
		).
		Where(
			goqu.I("asset.user_id").Eq(req.UserId),
		)

	if req.Keyword != "" {
		dataset = dataset.Where(goqu.I("ac.name").ILike("%" + req.Keyword + "%"))
	}

	result = make([]model.GetAsset, 0)
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		countDataset := dataset.Select(goqu.COUNT("*").As("total"))

		countSQL, countVals, err := countDataset.ToSQL()
		if err != nil {
			return fmt.Errorf("failed to build count SQL: %w", err)
		}

		if err := db.Get(&total, countSQL, countVals...); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to query count: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		dataset := libs.PaginationRequest(dataset, req.PaginationRequest)

		sql, val, err := dataset.ToSQL()
		if err != nil {
			return fmt.Errorf("failed to build SQL query: %w", err)
		}

		row, err := db.Queryx(sql, val...)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
		defer row.Close()

		err = libs.ScanRowsIntoStructs(row, &result)
		if err != nil {
			return fmt.Errorf("failed to scan rows into structs: %w", err)
		}

		return nil
	})

	err = g.Wait()
	if err != nil {
		return nil, 0, err
	}

	return
}

func (r *repository) Update(userId, id uint, data map[string]any, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	selectQ, selectV, err := dialect.From("asset").
		Where(
			goqu.I("id").Eq(id),
			goqu.I("user_id").Eq(userId),
		).
		ForUpdate(exp.Wait).
		ToSQL()

	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	_, err = tx.Exec(selectQ, selectV...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	dataset := dialect.Update("asset").Set(data).
		Where(
			goqu.I("id").Eq(id),
			goqu.I("user_id").Eq(userId),
		)

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

func (r *repository) Delete(userId, id uint, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.Delete("asset").
		Where(
			goqu.I("user_id").Eq(userId),
			goqu.I("id").Eq(id),
		)

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
