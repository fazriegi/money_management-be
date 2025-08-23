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

type IAssetUsecase interface {
	GetAssets(req *model.AssetRequest) (resp model.Response)
	Update(req *model.InsertAssetRequest) (resp model.Response)
}

type AssetUsecase struct {
	repository repository.IAssetRepository
	log        *logrus.Logger
}

func NewAssetUsecase(repo repository.IAssetRepository) IAssetUsecase {
	log := config.GetLogger()

	return &AssetUsecase{
		repository: repo,
		log:        log,
	}
}

func (u *AssetUsecase) GetAssets(req *model.AssetRequest) (resp model.Response) {
	db := config.GetDatabase()

	listData, err := u.repository.GetAssets(req, db)
	if err != nil {
		u.log.Errorf("repository.GetAssets: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	respData := make([]model.Asset, len(listData))
	for i, data := range listData {
		decAmount, err := libs.Decrypt(data.Amount.(string))
		if err != nil {
			u.log.Errorf("error decrypting amount: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		amount, _ := strconv.Atoi(decAmount)

		decValue, err := libs.Decrypt(data.Value.(string))
		if err != nil {
			u.log.Errorf("error decrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		value, _ := strconv.Atoi(decValue)

		respData[i] = model.Asset{
			PeriodCode: data.PeriodCode,
			Name:       data.Name,
			Amount:     amount,
			Value:      value,
			OrderNo:    data.OrderNo,
		}
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")
	resp.Data = respData
	return
}

func (u *AssetUsecase) Update(req *model.InsertAssetRequest) (resp model.Response) {
	db := config.GetDatabase()
	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error begin tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	insertData := make([]model.Asset, 0)
	for _, data := range req.Data {
		encAmount, err := libs.Encrypt(fmt.Sprintf("%0.f", data.Amount))
		if err != nil {
			u.log.Errorf("error encrypting amount: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		encValue, err := libs.Encrypt(fmt.Sprintf("%0.f", data.Value))
		if err != nil {
			u.log.Errorf("error encrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		insertData = append(insertData, model.Asset{
			PeriodCode: req.PeriodCode,
			Name:       data.Name,
			Amount:     encAmount,
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
