package cashflow

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/cashflow/model"
	"github.com/fazriegi/money_management-be/module/common"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/sirupsen/logrus"
)

type Usecase interface {
	List(user *userModel.User, req *model.ListRequest) (resp common.Response)
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

func (u *usecase) List(user *userModel.User, req *model.ListRequest) (resp common.Response) {
	db := config.GetDatabase()

	req.UserId = user.ID
	listData, total, err := u.repo.List(req, db)
	if err != nil {
		u.log.Errorf("repo.List: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	result := make([]model.CashflowData, len(listData))
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

		result[i] = model.CashflowData{
			ID:       data.ID,
			Category: data.Category,
			Date:     data.Date,
			Value:    value,
			Type:     data.Type,
		}
	}

	responseData := map[string]any{
		"data":  result,
		"total": total,
	}

	return resp.CustomResponse(http.StatusOK, "success", responseData)
}
