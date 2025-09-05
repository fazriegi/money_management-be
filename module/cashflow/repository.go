package cashflow

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/cashflow/expense"
	"github.com/fazriegi/money_management-be/module/cashflow/income"
	"github.com/fazriegi/money_management-be/module/cashflow/model"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Repository interface {
	List(req *model.ListRequest, db *sqlx.DB) (result []model.GetCashflow, total uint, err error)
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

func (r *repository) List(req *model.ListRequest, db *sqlx.DB) (result []model.GetCashflow, total uint, err error) {
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
	result = make([]model.GetCashflow, 0)
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		countDataset := dataset.Select(goqu.COUNT("*").As("total"))

		countSQL, countVals, err := countDataset.ToSQL()
		if err != nil {
			return fmt.Errorf("failed to build count SQL: %w", err)
		}

		if err := db.Get(&total, countSQL, countVals...); err != nil {
			return fmt.Errorf("failed to query count: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		dataset = libs.PaginationRequest(dataset, req.PaginationRequest)

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
