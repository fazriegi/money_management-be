package usecase

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/sirupsen/logrus"
)

type ILiabilityUsecase interface {
	GetList(user *model.User, req *model.GetLiabilityRequest) (resp model.Response)
	Update(user *model.User, req *model.UpdateLiabilityRequest) (resp model.Response)
}

type LiabilityUsecase struct {
	repository repository.ILiabilityRepository
	log        *logrus.Logger
}

func NewLiabilityUsecase(repo repository.ILiabilityRepository) ILiabilityUsecase {
	log := config.GetLogger()

	return &LiabilityUsecase{
		repository: repo,
		log:        log,
	}
}

func (u *LiabilityUsecase) GetList(user *model.User, req *model.GetLiabilityRequest) (resp model.Response) {
	db := config.GetDatabase()

	listData, err := u.repository.GetList(req, user.ID, db)
	if err != nil {
		u.log.Errorf("repository.GetList: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	respData := make(model.GetLiabilityResponse, len(listData))
	for i, data := range listData {
		decValue, err := libs.Decrypt(data.Value.(string))
		if err != nil {
			u.log.Errorf("error decrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
			return
		}

		value, _ := strconv.Atoi(decValue)

		respData[i] = model.LiabilityResponse{
			ID:         data.ID,
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

func (u *LiabilityUsecase) Update(user *model.User, req *model.UpdateLiabilityRequest) (resp model.Response) {
	db := config.GetDatabase()
	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error begin tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}
	defer tx.Rollback()

	keepID := []uint{}
	insertData := make([]model.Liability, 0)
	for _, data := range req.Data {
		encValue, err := libs.Encrypt(fmt.Sprintf("%0.f", data.Value))
		if err != nil {
			u.log.Errorf("error encrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
			return
		}

		if data.ID != 0 {
			keepID = append(keepID, data.ID)
			err = u.repository.UpdateByID(tx, data.ID, user.ID, map[string]any{
				"period_code": req.PeriodCode,
				"name":        data.Name,
				"value":       encValue,
				"order_no":    data.OrderNo,
			})
			if err != nil {
				u.log.Errorf("repository.UpdateByID: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
				return
			}
		} else {
			insertData = append(insertData, model.Liability{
				PeriodCode: req.PeriodCode,
				Name:       data.Name,
				Value:      encValue,
				OrderNo:    data.OrderNo,
				UserID:     user.ID,
			})
		}
	}

	err = u.repository.DeleteExcept(tx, keepID, req.PeriodCode, user.ID)
	if err != nil {
		u.log.Errorf("repository.DeleteExcept: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	if len(insertData) > 0 {
		err = u.repository.BulkInsert(tx, &insertData)
		if err != nil {
			u.log.Errorf("repository.BulkInsert: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		u.log.Errorf("error committing tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")

	return
}
