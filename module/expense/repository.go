package expense

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/expense/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Insert(data *model.Expense, tx *sqlx.Tx) error
	List(req *model.ListRequest, db *sqlx.DB) (result []model.Expense, err error)
	ListCategory(userID uint, db *sqlx.DB) (result []model.ExpenseCategory, err error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Insert(data *model.Expense, tx *sqlx.Tx) error {
	dialect := libs.GetDialect()

	dataset := dialect.Insert("expense").Rows(*data)
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

func (r *repository) List(req *model.ListRequest, db *sqlx.DB) (result []model.Expense, err error) {
	dialect := libs.GetDialect()

	if req.Sort == nil {
		sort := "date desc"
		req.Sort = &sort
	}

	dataset := dialect.
		From("expense").
		Join(goqu.T("user_expense_category").As("uec"), goqu.On(
			goqu.I("uec.id").Eq(goqu.I("expense.category_id")),
			goqu.I("uec.user_id").Eq(goqu.I("expense.user_id")),
		)).
		Select(
			goqu.I("expense.id"),
			goqu.I("expense.category_id"),
			goqu.I("uec.name").As("category"),
			goqu.I("expense.date"),
			goqu.I("expense.value"),
			goqu.I("expense.notes"),
		).
		Where(
			goqu.I("expense.user_id").Eq(req.UserId),
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

	result = make([]model.Expense, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}

func (r *repository) ListCategory(userID uint, db *sqlx.DB) (result []model.ExpenseCategory, err error) {
	dialect := libs.GetDialect()

	dataset := dialect.From("user_expense_category").Where(goqu.I("user_id").Eq(userID))

	sql, val, err := dataset.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	row, err := db.Queryx(sql, val...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer row.Close()

	result = make([]model.ExpenseCategory, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}
