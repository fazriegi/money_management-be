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

type IUsecase interface {
	GetAssets(req model.AssetRequest) (resp model.Response)
	Update(req *model.InsertAssetRequest) (resp model.Response)
}

type Usecase struct {
	repository *repository.Repository
	log        *logrus.Logger
}

func NewUsecase(repo *repository.Repository) IUsecase {
	log := config.GetLogger()

	return &Usecase{
		repository: repo,
		log:        log,
	}
}

func (u *Usecase) GetAssets(req model.AssetRequest) (resp model.Response) {
	assets, err := u.repository.GetAssets(req)
	if err != nil {
		u.log.Errorf("repository.GetAssets: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	respData := make([]model.Asset, len(assets))
	for i, asset := range assets {
		decAmount, err := libs.Decrypt(asset.Amount.(string))
		if err != nil {
			u.log.Errorf("error decrypting amount: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		amount, _ := strconv.Atoi(decAmount)

		decValue, err := libs.Decrypt(asset.Value.(string))
		if err != nil {
			u.log.Errorf("error decrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		value, _ := strconv.Atoi(decValue)

		respData[i] = model.Asset{
			PeriodCode: asset.PeriodCode,
			Name:       asset.Name,
			Amount:     amount,
			Value:      value,
			OrderNo:    asset.OrderNo,
		}
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")
	resp.Data = respData
	return
}

func (u *Usecase) Update(req *model.InsertAssetRequest) (resp model.Response) {
	db := config.GetDatabase()
	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error begin tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	insertData := make([]model.Asset, 0)
	for _, asset := range req.Data {
		encAmount, err := libs.Encrypt(fmt.Sprintf("%0.f", asset.Amount))
		if err != nil {
			u.log.Errorf("error encrypting amount: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		encValue, err := libs.Encrypt(fmt.Sprintf("%0.f", asset.Value))
		if err != nil {
			u.log.Errorf("error encrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
			return
		}

		insertData = append(insertData, model.Asset{
			PeriodCode: req.PeriodCode,
			Name:       asset.Name,
			Amount:     encAmount,
			Value:      encValue,
			OrderNo:    asset.OrderNo,
		})
	}

	err = u.repository.DeleteAssetByPeriod(tx, req.PeriodCode)
	if err != nil {
		u.log.Errorf("repository.DeleteAssetByPeriod: %s", err.Error())
		tx.Rollback()
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	err = u.repository.BulkInsert(tx, insertData)
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
