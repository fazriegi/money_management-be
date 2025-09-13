package cashflow

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/cashflow/expense"
	expenseModel "github.com/fazriegi/money_management-be/module/cashflow/expense/model"
	"github.com/fazriegi/money_management-be/module/cashflow/income"
	incomeModel "github.com/fazriegi/money_management-be/module/cashflow/income/model"
	"github.com/fazriegi/money_management-be/module/cashflow/model"
	"github.com/fazriegi/money_management-be/module/common"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Usecase interface {
	List(user *userModel.User, req *model.ListRequest) (resp common.Response)
}

type usecase struct {
	log            *logrus.Logger
	repo           Repository
	incomeUsecase  income.Usecase
	expenseUsecase expense.Usecase
}

func NewUsecase(log *logrus.Logger, repo Repository, incomeUsecase income.Usecase, expenseUsecase expense.Usecase) Usecase {
	return &usecase{
		log,
		repo,
		incomeUsecase,
		expenseUsecase,
	}
}

func (u *usecase) List(user *userModel.User, req *model.ListRequest) (resp common.Response) {
	db := config.GetDatabase()

	var (
		totalData    uint
		result       []model.CashflowData
		totalIncome  float64
		totalExpense float64
	)

	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {

		req.UserId = user.ID
		listData, total, err := u.repo.List(req, db)
		if err != nil {
			u.log.Errorf("repo.List: %s", err.Error())
			return errors.New("failed list data")
		}

		totalData = total

		resultData := make([]model.CashflowData, len(listData))
		for i, data := range listData {
			decValue, err := libs.Decrypt(fmt.Sprintf("%d", user.ID), data.Value)
			if err != nil {
				u.log.Errorf("error decrypting value: %s", err.Error())
				return errors.New("failed list data")
			}

			value, err := strconv.ParseFloat(decValue, 64)
			if err != nil {
				u.log.Errorf("error parsing string: %s", err.Error())
				return errors.New("failed list data")
			}

			resultData[i] = model.CashflowData{
				ID:       data.ID,
				Category: data.Category,
				Date:     data.Date,
				Value:    value,
				Type:     data.Type,
			}
		}

		result = resultData

		return nil
	})

	g.Go(func() error {
		resp = u.incomeUsecase.List(user, &incomeModel.ListRequest{
			UserId:    user.ID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
		})
		if !resp.IsSuccess {
			return errors.New("failed calculate total income")
		}

		data := resp.Data.([]incomeModel.IncomeData)

		for _, v := range data {
			totalIncome += v.Value
		}

		return nil
	})

	g.Go(func() error {
		resp = u.expenseUsecase.List(user, &expenseModel.ListRequest{
			UserId:    user.ID,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
		})
		if !resp.IsSuccess {
			return errors.New("failed calculate total expense")
		}

		data := resp.Data.([]expenseModel.ExpenseData)

		for _, v := range data {
			totalExpense += v.Value
		}

		return nil
	})

	err := g.Wait()
	if err != nil {
		u.log.Errorf("%s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	responseData := map[string]any{
		"cashflow": map[string]any{
			"data":           result,
			"total_income":   totalIncome,
			"total_expense":  totalExpense,
			"total_cashflow": totalIncome - totalExpense,
		},
		"total": totalData,
	}

	return resp.CustomResponse(http.StatusOK, "success", responseData)
}
