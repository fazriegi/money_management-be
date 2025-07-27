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

type ILiabilityUsecase interface {
	GetList(req *model.GetLiabilityRequest) (resp model.Response)
	Update(req *model.UpdateLiabilityRequest) (resp model.Response)
}

type LiabilityUsecase struct {
	repository *repository.LiabilityRepository
	log        *logrus.Logger
}

func NewLiabilityUsecase(repo *repository.LiabilityRepository) ILiabilityUsecase {
	log := config.GetLogger()

	return &LiabilityUsecase{
		repository: repo,
		log:        log,
	}
}

func (u *LiabilityUsecase) GetList(req *model.GetLiabilityRequest) (resp model.Response) {
	listData, err := u.repository.GetList(req)
	if err != nil {
		u.log.Errorf("repository.GetList: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	respData := make(model.GetLiabilityResponse, len(listData))
	for i, data := range listData {
		decValue, err := libs.Decrypt(data.Value.(string))
		if err != nil {
			u.log.Errorf("error decrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		value, _ := strconv.Atoi(decValue)

		respData[i] = model.LiabilityResponse{
			PeriodCode: data.PeriodCode,
			Name:       data.Name,
			Value:      value,
			OrderNo:    data.OrderNo,
		}
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")
	resp.Data = respData
	return
}

func (u *LiabilityUsecase) Update(req *model.UpdateLiabilityRequest) (resp model.Response) {
	db := config.GetDatabase()
	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error begin tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	insertData := make([]model.Liability, 0)
	for _, data := range req.Data {
		encValue, err := libs.Encrypt(fmt.Sprintf("%0.f", data.Value))
		if err != nil {
			u.log.Errorf("error encrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		insertData = append(insertData, model.Liability{
			PeriodCode: req.PeriodCode,
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
