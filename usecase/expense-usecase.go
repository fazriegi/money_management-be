package usecase

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/sirupsen/logrus"
)

type IExpenseUsecase interface {
	GetList(req *model.GetExpenseRequest) (resp model.Response)
	Update(req *model.UpdateExpenseRequest) (resp model.Response)
}

type ExpenseUsecase struct {
	repository repository.IExpenseRepository
	log        *logrus.Logger
}

func NewExpenseUsecase(repo repository.IExpenseRepository) IExpenseUsecase {
	log := config.GetLogger()

	return &ExpenseUsecase{
		repository: repo,
		log:        log,
	}
}

func (u *ExpenseUsecase) GetList(req *model.GetExpenseRequest) (resp model.Response) {
	db := config.GetDatabase()

	listData, err := u.repository.GetList(req, db)
	if err != nil {
		u.log.Errorf("repository.GetList: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	respData := make(model.GetExpenseResponse, len(listData))
	for i, data := range listData {
		decValue, err := libs.Decrypt(data.Value.(string))
		if err != nil {
			u.log.Errorf("error decrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		value, _ := strconv.Atoi(decValue)

		respData[i] = model.ExpenseResponse{
			PeriodCode:  data.PeriodCode,
			Name:        data.Name,
			Value:       value,
			OrderNo:     data.OrderNo,
			LiabilityID: data.LiabilityID,
		}
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")
	resp.Data = respData
	return
}

func (u *ExpenseUsecase) Update(req *model.UpdateExpenseRequest) (resp model.Response) {
	db := config.GetDatabase()
	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error begin tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	insertData := make([]model.Expense, 0)
	for _, data := range req.Data {
		encValue, err := libs.Encrypt(fmt.Sprintf("%0.f", data.Value))
		if err != nil {
			u.log.Errorf("error encrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		insertData = append(insertData, model.Expense{
			PeriodCode:  req.PeriodCode,
			Name:        data.Name,
			Value:       encValue,
			OrderNo:     data.OrderNo,
			LiabilityID: data.LiabilityID,
		})
	}

	err = u.repository.DeleteByPeriod(tx, req.PeriodCode)
	if err != nil {
		u.log.Errorf("repository.DeleteByPeriod: %s", err.Error())
		tx.Rollback()
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	if len(insertData) > 0 {
		err = u.repository.BulkInsert(tx, &insertData)
		if err != nil {
			u.log.Errorf("repository.BulkInsert: %s", err.Error())
			tx.Rollback()
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		u.log.Errorf("error committing tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")

	return
}
