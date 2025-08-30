package repository

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/jmoiron/sqlx"
)

type ILiabilityRepository interface {
	GetList(req *model.GetLiabilityRequest, userID uint, db *sqlx.DB) (result []model.Liability, err error)
	BulkInsert(tx *sqlx.Tx, data *[]model.Liability) error
	DeleteByPeriod(tx *sqlx.Tx, periodCode string, userID uint) error
	DeleteExcept(tx *sqlx.Tx, keepId []uint, periodCode string, userID uint) error
	UpdateByID(tx *sqlx.Tx, id, userID uint, data map[string]any) error
	GetByID(id, userID uint, db *sqlx.DB) (result model.Liability, err error)
}

type LiabilityRepository struct {
}

func NewLiabilityRepository() ILiabilityRepository {
	return &LiabilityRepository{}
}

func (r *LiabilityRepository) GetList(req *model.GetLiabilityRequest, userID uint, db *sqlx.DB) (result []model.Liability, err error) {
	dialect := libs.GetDialect()

	if req.Sort == nil || *req.Sort == "" {
		order := "order_no" // Default sorting if not provided
		req.Sort = &order
	}

	dataset := dialect.From("liabilities").Where(goqu.I("user_id").Eq(userID))

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

func (r *LiabilityRepository) DeleteByPeriod(tx *sqlx.Tx, periodCode string, userID uint) error {
	dialect := libs.GetDialect()

	dataset := dialect.Delete("liabilities").Where(goqu.I("period_code").Eq(periodCode), goqu.I("user_id").Eq(userID))
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

func (r *LiabilityRepository) GetByID(id, userID uint, db *sqlx.DB) (result model.Liability, err error) {
	dialect := libs.GetDialect()

	dataset := dialect.From("liabilities").
		Select(
			goqu.I("id"),
			goqu.I("period_code"),
			goqu.I("name"),
			goqu.I("value"),
			goqu.I("order_no"),
		).
		Where(
			goqu.I("id").Eq(id),
			goqu.I("user_id").Eq(userID),
		)

	query, val, err := dataset.ToSQL()
	if err != nil {
		return result, fmt.Errorf("failed to build SQL query: %w", err)
	}

	err = db.Get(&result, query, val...)
	if err != nil {
		return result, err
	}

	return
}

func (r *LiabilityRepository) UpdateByID(tx *sqlx.Tx, id, userID uint, data map[string]any) error {
	dialect := libs.GetDialect()
	selectQ, selectV, err := dialect.From("liabilities").
		Where(
			goqu.I("id").Eq(id),
			goqu.I("user_id").Eq(userID),
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

	dataset := dialect.Update("liabilities").Set(data).Where(
		goqu.I("id").Eq(id),
		goqu.I("user_id").Eq(userID),
	)

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

func (r *LiabilityRepository) DeleteExcept(tx *sqlx.Tx, keepId []uint, periodCode string, userID uint) error {
	dialect := libs.GetDialect()

	dataset := dialect.Delete("liabilities").
		Where(
			goqu.I("period_code").Eq(periodCode),
			goqu.I("user_id").Eq(userID),
		)

	if len(keepId) > 0 {
		dataset = dataset.Where(goqu.I("id").NotIn(keepId))
	}

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
