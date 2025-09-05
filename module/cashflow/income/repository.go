package income

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/cashflow/income/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Insert(data *model.Income, tx *sqlx.Tx) error
	ListCategory(userID uint, db *sqlx.DB) (result []model.IncomeCategory, err error)
	List(req *model.ListRequest, db *sqlx.DB) (result []model.GetIncome, err error)
	Update(userId, id uint, data map[string]any, tx *sqlx.Tx) error
	Delete(userId, id uint, tx *sqlx.Tx) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Insert(data *model.Income, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("income").Rows(*data)
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

func (r *repository) ListCategory(userID uint, db *sqlx.DB) (result []model.IncomeCategory, err error) {
	dialect := libs.GetDialect()

	dataset := dialect.From("user_income_category").Where(goqu.I("user_id").Eq(userID))

	sql, val, err := dataset.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	row, err := db.Queryx(sql, val...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer row.Close()

	result = make([]model.IncomeCategory, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *repository) List(req *model.ListRequest, db *sqlx.DB) (result []model.GetIncome, err error) {
	dialect := libs.GetDialect()

	if req.Sort == nil {
		sort := "date desc"
		req.Sort = &sort
	}

	dataset := dialect.
		From("income").
		Join(goqu.T("user_income_category").As("uec"), goqu.On(
			goqu.I("uec.id").Eq(goqu.I("income.category_id")),
			goqu.I("uec.user_id").Eq(goqu.I("income.user_id")),
		)).
		Select(
			goqu.I("income.id"),
			goqu.I("income.category_id"),
			goqu.I("uec.name").As("category"),
			goqu.I("income.date"),
			goqu.I("income.value"),
			goqu.I("income.notes"),
		).
		Where(
			goqu.I("income.user_id").Eq(req.UserId),
		)

	if req.Keyword != "" {
		dataset = dataset.Where(goqu.I("uec.name").Eq(req.Keyword))
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

	result = make([]model.GetIncome, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *repository) Update(userId, id uint, data map[string]any, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	selectQ, selectV, err := dialect.From("income").
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

	dataset := dialect.Update("income").Set(data).
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

	dataset := dialect.Delete("income").
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
