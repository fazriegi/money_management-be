package expense

import (
	"fmt"
	"net/http"

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

func (u *usecase) ListCategory(user *userModel.User) (resp common.Response) {
	db := config.GetDatabase()

	data, err := u.repo.ListCategory(user.ID, db)
	if err != nil {
		u.log.Errorf("repo.ListCategory: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	return resp.CustomResponse(http.StatusOK, "success", data)
}
