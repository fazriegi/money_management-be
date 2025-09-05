package cashflow

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/cashflow/expense"
	"github.com/fazriegi/money_management-be/module/cashflow/income"
	"github.com/fazriegi/money_management-be/module/cashflow/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	List(req *model.ListRequest, db *sqlx.DB) (result []model.GetCashflow, err error)
}

type repository struct {
	expenseRepo expense.Repository
	incomeRepo  income.Repository
}

func NewRepository(expenseRepo expense.Repository, incomeRepo income.Repository) Repository {
	return &repository{
		expenseRepo,
		incomeRepo,
	}
}

func (r *repository) List(req *model.ListRequest, db *sqlx.DB) (result []model.GetCashflow, err error) {
	dialect := libs.GetDialect()

	if req.Sort == nil {
		sort := "date desc"
		req.Sort = &sort
	}

	listFilter := model.ListFilter{
		UserId:    req.UserId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	expenseDataset := r.expenseRepo.CreateListQuery(&listFilter)
	incomeDataset := r.incomeRepo.CreateListQuery(&listFilter)
	unionDataset := expenseDataset.Union(incomeDataset)

	dataset := dialect.
		From(unionDataset.As("obj")).
		Select(
			goqu.I("id"),
			goqu.I("category"),
			goqu.I("date"),
			goqu.I("value"),
			goqu.I("user_id"),
		)

	if req.Category != "" {
		dataset = dataset.Where(goqu.I("category").ILike("%" + req.Category + "%"))
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

	result = make([]model.GetCashflow, 0)
	err = libs.ScanRowsIntoStructs(row, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to scan rows into structs: %w", err)
	}

	return
}
