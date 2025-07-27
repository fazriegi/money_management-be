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

type IIncomeUsecase interface {
	GetList(req *model.GetIncomeRequest) (resp model.Response)
	Update(req *model.UpdateIncomeRequest) (resp model.Response)
}

type IncomeUsecase struct {
	repository *repository.IncomeRepository
	log        *logrus.Logger
}

func NewIncomeUsecase(repo *repository.IncomeRepository) IIncomeUsecase {
	log := config.GetLogger()

	return &IncomeUsecase{
		repository: repo,
		log:        log,
	}
}

func (u *IncomeUsecase) GetList(req *model.GetIncomeRequest) (resp model.Response) {
	listData, err := u.repository.GetList(req)
	if err != nil {
		u.log.Errorf("repository.GetList: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	respData := make(model.GetIncomeResponse, len(listData))
	for i, data := range listData {
		decValue, err := libs.Decrypt(data.Value.(string))
		if err != nil {
			u.log.Errorf("error decrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		value, _ := strconv.Atoi(decValue)

		respData[i] = model.IncomeResponse{
			PeriodCode: data.PeriodCode,
			Type:       data.Type,
			Name:       data.Name,
			Value:      value,
			OrderNo:    data.OrderNo,
		}
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")
	resp.Data = respData
	return
}

func (u *IncomeUsecase) Update(req *model.UpdateIncomeRequest) (resp model.Response) {
	db := config.GetDatabase()
	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error begin tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	insertData := make([]model.Income, 0)
	for _, data := range req.Data {
		encValue, err := libs.Encrypt(fmt.Sprintf("%0.f", data.Value))
		if err != nil {
			u.log.Errorf("error encrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		insertData = append(insertData, model.Income{
			PeriodCode: req.PeriodCode,
			Type:       data.Type,
			Name:       data.Name,
			Value:      encValue,
			OrderNo:    data.OrderNo,
		})
	}

	err = u.repository.DeleteByPeriod(tx, req.PeriodCode)
	if err != nil {
		u.log.Errorf("repository.DeleteByPeriod: %s", err.Error())
		tx.Rollback()
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	err = u.repository.BulkInsert(tx, &insertData)
	if err != nil {
		u.log.Errorf("repository.BulkInsert: %s", err.Error())
		tx.Rollback()
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	if err = tx.Commit(); err != nil {
		u.log.Errorf("error committing tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")

	return
}
