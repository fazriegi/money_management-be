package user

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetByUsername(username string, db *sqlx.DB) (model.User, error)
	Insert(data *model.User, db *sqlx.Tx) (result uint, err error)
	CreateIncomeCat(userId uint, tx *sqlx.Tx) error
	CreateExpenseCat(userId uint, tx *sqlx.Tx) error
	CreateAssetCat(userId uint, tx *sqlx.Tx) error
	CreatePeriod(userId uint, tx *sqlx.Tx) error
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetByUsername(username string, db *sqlx.DB) (result model.User, err error) {
	dialect := libs.GetDialect()

	dataset := dialect.From("user").Select(goqu.I("username"), goqu.I("password"), goqu.I("email"), goqu.I("id"), goqu.I("name")).Where(goqu.I("username").Eq(username))

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

func (r *repository) Insert(data *model.User, tx *sqlx.Tx) (result uint, err error) {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("user").Rows(*data)
	sql, val, err := dataset.ToSQL()
	if err != nil {
		return result, fmt.Errorf("failed to build SQL query: %w", err)
	}

	res, err := tx.Exec(sql, val...)
	if err != nil {
		return result, fmt.Errorf("failed to execute insert: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return result, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return uint(id), nil
}

func (r *repository) CreateIncomeCat(userId uint, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.
		Insert("user_expense_category").
		Cols(
			goqu.I("name"),
			goqu.I("user_id"),
		).
		FromQuery(
			goqu.From("expense_category_default").
				Select(
					goqu.I("name"),
					goqu.L("?", userId),
				),
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

func (r *repository) CreateExpenseCat(userId uint, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.
		Insert("user_income_category").
		Cols(
			goqu.I("name"),
			goqu.I("user_id"),
		).
		FromQuery(
			goqu.From("income_category_default").
				Select(
					goqu.I("name"),
					goqu.L("?", userId),
				),
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

func (r *repository) CreateAssetCat(userId uint, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.
		Insert("asset_category").
		Cols(
			goqu.I("name"),
			goqu.I("user_id"),
		).
		FromQuery(
			goqu.From("asset_category_default").
				Select(
					goqu.I("name"),
					goqu.L("?", userId),
				),
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

func (r *repository) CreatePeriod(userId uint, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.
		Insert("monthly_period").
		Rows(
			map[string]interface{}{"day_of_month": 1, "user_id": userId},
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
