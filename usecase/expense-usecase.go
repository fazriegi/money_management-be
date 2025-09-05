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
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	respData := make(model.GetExpenseResponse, len(listData))
	for i, data := range listData {
		decValue, err := libs.Decrypt("", data.Value.(string))
		if err != nil {
			u.log.Errorf("error decrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
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
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}
	defer tx.Rollback()

	refetchLiability := false

	existingData, err := u.repository.GetListForUpdate(req.PeriodCode, user.ID, tx)
	if err != nil {
		u.log.Errorf("error encrypting value: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	liabilityMap := make(map[uint]float64)
	// return liability value
	for _, curData := range existingData {
		if curData.LiabilityID != nil {
			expensesDecValue, err := libs.Decrypt("", string(curData.Value.(string)))
			if err != nil {
				u.log.Errorf("error decrypting value: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
				return
			}

			expensesValue, err := strconv.ParseFloat(expensesDecValue, 64)
			if err != nil {
				u.log.Errorf("error convert nominal: %s", err.Error())
				resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
				return
			}

			if _, ok := liabilityMap[*curData.LiabilityID]; ok {
				liabilityMap[*curData.LiabilityID] += expensesValue
			} else {
				liability, err := u.liabilityRepo.GetByID(*curData.LiabilityID, user.ID, db)
				if err != nil {
					u.log.Errorf("liabilityRepo.GetByID: %s", err.Error())
					resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
					return
				}

				liabilityDecValue, err := libs.Decrypt("", string(liability.Value))
				if err != nil {
					u.log.Errorf("error decrypting value: %s", err.Error())
					resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
					return
				}

				liabilityValue, err := strconv.ParseFloat(liabilityDecValue, 64)
				if err != nil {
					u.log.Errorf("error convert nominal: %s", err.Error())
					resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
					return
				}

				liabilityNewValue := liabilityValue + expensesValue
				liabilityMap[*curData.LiabilityID] = liabilityNewValue
			}
		}
	}

	err = u.repository.DeleteByPeriod(tx, req.PeriodCode, user.ID)
	if err != nil {
		u.log.Errorf("repository.DeleteByPeriod: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	insertData := make([]model.Expense, 0)
	for _, data := range req.Data {
		encValue, err := libs.Encrypt("", fmt.Sprintf("%0.f", data.Value))
		if err != nil {
			u.log.Errorf("error encrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
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
			if _, ok := liabilityMap[*data.LiabilityID]; ok {
				liabilityMap[*data.LiabilityID] -= data.Value.(float64)
			} else {
				liability, err := u.liabilityRepo.GetByID(*data.LiabilityID, user.ID, db)
				if err != nil {
					u.log.Errorf("liabilityRepo.GetByID: %s", err.Error())
					resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
					return
				}

				liabilityDecValue, err := libs.Decrypt("", string(liability.Value))
				if err != nil {
					u.log.Errorf("error decrypting value: %s", err.Error())
					resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
					return
				}

				liabilityValue, err := strconv.ParseFloat(liabilityDecValue, 64)
				if err != nil {
					u.log.Errorf("error convert nominal: %s", err.Error())
					resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
					return
				}

				liabilityNewValue := liabilityValue - data.Value.(float64)
				liabilityMap[*data.LiabilityID] = liabilityNewValue
			}
		}
	}

	for id, val := range liabilityMap {
		refetchLiability = true

		encValue, err := libs.Encrypt("", fmt.Sprintf("%0.f", val))
		if err != nil {
			u.log.Errorf("error encrypting value: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
			return
		}

		err = u.liabilityRepo.UpdateByID(tx, id, user.ID, map[string]any{"value": encValue})
		if err != nil {
			u.log.Errorf("liabilityRepo.UpdateByID: %s", err.Error())
			resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
			return
		}
	}

	err = u.repository.DeleteByPeriod(tx, req.PeriodCode, user.ID)
	if err != nil {
		u.log.Errorf("repository.DeleteByPeriod: %s", err.Error())
		tx.Rollback()
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	if len(insertData) > 0 {
		err = u.repository.BulkInsert(tx, &insertData)
		if err != nil {
			u.log.Errorf("repository.BulkInsert: %s", err.Error())
			tx.Rollback()
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
	resp.Data = map[string]any{"refetch_liability": refetchLiability}

	return
}
