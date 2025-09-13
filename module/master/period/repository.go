package period

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/master/period/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetPeriod(userId uint, db *sqlx.DB) (result model.MonthlyPeriod, err error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetPeriod(userId uint, db *sqlx.DB) (result model.MonthlyPeriod, err error) {
	dialect := libs.GetDialect()

	dataset := dialect.From("monthly_period").
		Select(goqu.I("day_of_month")).
		Where(goqu.I("user_id").Eq(userId))

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
