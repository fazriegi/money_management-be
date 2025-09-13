package period

import (
	"net/http"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/module/common"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/sirupsen/logrus"
)

type Usecase interface {
	Get(user *userModel.User) (resp common.Response)
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

func (u *usecase) Get(user *userModel.User) (resp common.Response) {
	db := config.GetDatabase()

	result, err := u.repo.GetPeriod(user.ID, db)
	if err != nil {
		u.log.Errorf("repo.List: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	return resp.CustomResponse(http.StatusOK, "success", result)
}
