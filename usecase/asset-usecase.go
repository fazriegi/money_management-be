package usecase

import (
	"net/http"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/fazriegi/money_management-be/repository"
)

type IUsecase interface {
	GetAssets(req model.AssetRequest) (resp model.Response)
}

type Usecase struct {
	repository *repository.Repository
}

func NewUsecase(repo *repository.Repository) IUsecase {
	return &Usecase{
		repository: repo,
	}
}

func (u *Usecase) GetAssets(req model.AssetRequest) (resp model.Response) {
	var (
		log = config.GetLogger()
	)

	assets, err := u.repository.GetAssets(req)
	if err != nil {
		log.Errorf("repository.GetAssets: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, "unexpected error occured")
		return
	}

	resp.Status = libs.CustomResponse(http.StatusOK, "success")
	resp.Data = assets
	return
}
