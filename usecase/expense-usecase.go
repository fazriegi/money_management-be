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
	GetList(user *model.User, req *model.GetExpenseRequest) (resp model.Response)
	Update(user *model.User, req *model.UpdateExpenseRequest) (resp model.Response)
}

type ExpenseUsecase struct {
	repository    repository.IExpenseRepository
	liabilityRepo repository.ILiabilityRepository
	log           *logrus.Logger
}

func NewExpenseUsecase(repo repository.IExpenseRepository, liabilityRepo repository.ILiabilityRepository) IExpenseUsecase {
	log := config.GetLogger()

	return &ExpenseUsecase{
		repository:    repo,
		log:           log,
		liabilityRepo: liabilityRepo,
	}
}

func (u *ExpenseUsecase) GetList(user *model.User, req *model.GetExpenseRequest) (resp model.Response) {
	db := config.GetDatabase()

	listData, err := u.repository.GetList(req, user.ID, db)
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

func (u *ExpenseUsecase) Update(user *model.User, req *model.UpdateExpenseRequest) (resp model.Response) {
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
			UserID:      user.ID,
		})

		if data.LiabilityID != nil {
			liability, err := u.liabilityRepo.GetByID(*data.LiabilityID, user.ID, db)
			if err != nil {
				u.log.Errorf("liabilityRepo.GetByID: %s", err.Error())
				tx.Rollback()
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}

			oldValue := liability.Value.([]uint8)
			decValue, err := libs.Decrypt(string(oldValue))
			if err != nil {
				u.log.Errorf("error decrypting value: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}

			oldNominal, err := strconv.ParseFloat(decValue, 64)
			if err != nil {
				u.log.Errorf("error convert nominal: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}

			newValue := oldNominal - data.Value.(float64)

			encValue, err := libs.Encrypt(fmt.Sprintf("%0.f", newValue))
			if err != nil {
				u.log.Errorf("error encrypting value: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}

			err = u.liabilityRepo.UpdateByID(tx, *data.LiabilityID, user.ID, map[string]any{"value": encValue})
			if err != nil {
				u.log.Errorf("liabilityRepo.UpdateByID: %s", err.Error())
				tx.Rollback()
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}
		}
	}

	for _, deletedData := range req.Delete {
		if deletedData.LiabilityID != nil {
			liability, err := u.liabilityRepo.GetByID(*deletedData.LiabilityID, user.ID, db)
			if err != nil {
				u.log.Errorf("liabilityRepo.GetByID: %s", err.Error())
				tx.Rollback()
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}

			oldValue := liability.Value.([]uint8)
			decValue, err := libs.Decrypt(string(oldValue))
			if err != nil {
				u.log.Errorf("error decrypting value: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}

			oldNominal, err := strconv.ParseFloat(decValue, 64)
			if err != nil {
				u.log.Errorf("error convert nominal: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}

			newValue := oldNominal + deletedData.Value.(float64)

			encValue, err := libs.Encrypt(fmt.Sprintf("%0.f", newValue))
			if err != nil {
				u.log.Errorf("error encrypting value: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}

			err = u.liabilityRepo.UpdateByID(tx, *deletedData.LiabilityID, user.ID, map[string]any{"value": encValue})
			if err != nil {
				u.log.Errorf("liabilityRepo.UpdateByID: %s", err.Error())
				tx.Rollback()
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
				return
			}
		}
	}

	err = u.repository.DeleteByPeriod(tx, req.PeriodCode, user.ID)
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
