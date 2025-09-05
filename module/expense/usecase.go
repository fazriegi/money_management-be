package expense

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/common"
	"github.com/fazriegi/money_management-be/module/expense/model"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/sirupsen/logrus"
)

type Usecase interface {
	Add(user *userModel.User, req *model.AddRequest) (resp common.Response)
	List(user *userModel.User, req *model.ListRequest) (resp common.Response)
	ListCategory(user *userModel.User) (resp common.Response)
}

type usecase struct {
	log  *logrus.Logger
	repo Repository
}

func NewUsecase(log *logrus.Logger, repo Repository) Usecase {
	return &usecase{
		log,
		repo,
	}
}

func (u *usecase) Add(user *userModel.User, req *model.AddRequest) (resp common.Response) {
	db := config.GetDatabase()
	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error begin tx: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}
	defer tx.Rollback()

	encValue, err := libs.Encrypt(fmt.Sprintf("%d", user.ID), fmt.Sprintf("%0.f", req.Value))
	if err != nil {
		u.log.Errorf("error encrypting value: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	data := model.Expense{
		CategoryId: req.CategoryId,
		Date:       req.Date,
		Value:      encValue,
		UserId:     user.ID,
		Notes:      req.Notes,
	}

	err = u.repo.Insert(&data, tx)
	if err != nil {
		u.log.Errorf("failed insert expense: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	if err := tx.Commit(); err != nil {
		u.log.Errorf("failed commit tx: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	return resp.CustomResponse(http.StatusCreated, "success", nil)
}

func (u *usecase) List(user *userModel.User, req *model.ListRequest) (resp common.Response) {
	db := config.GetDatabase()

	req.UserId = user.ID
	listData, err := u.repo.List(req, db)
	if err != nil {
		u.log.Errorf("repo.List: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	result := make([]model.ExpenseData, len(listData))
	for i, data := range listData {
		decValue, err := libs.Decrypt(fmt.Sprintf("%d", user.ID), data.Value)
		if err != nil {
			u.log.Errorf("error decrypting value: %s", err.Error())
			return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
		}

		value, err := strconv.ParseFloat(decValue, 64)
		if err != nil {
			u.log.Errorf("error parsing string: %s", err.Error())
			return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
		}

		result[i] = model.ExpenseData{
			ID:         data.ID,
			CategoryId: data.CategoryId,
			Category:   data.Category,
			Date:       data.Date,
			Value:      value,
			Notes:      data.Notes,
		}
	}

	return resp.CustomResponse(http.StatusCreated, "success", result)
}

func (u *usecase) ListCategory(user *userModel.User) (resp common.Response) {
	db := config.GetDatabase()

	data, err := u.repo.ListCategory(user.ID, db)
	if err != nil {
		u.log.Errorf("repo.ListCategory: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	return resp.CustomResponse(http.StatusOK, "success", data)
}
