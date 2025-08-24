package usecase

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/fazriegi/money_management-be/repository"

	"github.com/sirupsen/logrus"
)

type IAuthUsecase interface {
	Register(props *model.RegisterRequest) (resp model.Response)
	Login(props *model.LoginRequest) (resp model.Response)
	// Me(ctx context.Context) (resp model.Response)
}

type AuthUsecase struct {
	repository repository.IUserRepository
	log        *logrus.Logger
	jwt        *libs.JWT
}

func NewAuthUsecase(repository repository.IUserRepository, jwt *libs.JWT) IAuthUsecase {
	log := config.GetLogger()

	return &AuthUsecase{
		repository,
		log,
		jwt,
	}
}

func (u *AuthUsecase) Register(props *model.RegisterRequest) (resp model.Response) {
	var (
		err            error
		user           model.User
		hashedPassword string
		db             = config.GetDatabase()
	)

	user, err = u.repository.GetUserByUsername(props.Username, db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.log.Errorf("repository.GetUserByUsername: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	if user.Username != "" {
		resp.Status = libs.CustomResponse(http.StatusBadRequest, "username already exists")
		return
	}

	if hashedPassword, err = libs.HashPassword(props.Password); err != nil {
		u.log.Errorf("libs.HashPassword: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	user = model.User{
		Name:     props.Name,
		Email:    props.Email,
		Password: hashedPassword,
		Username: props.Username,
	}

	if err := u.repository.InsertUser(&user, db); err != nil {
		u.log.Errorf("repository.InsertUser: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	createdUser := model.UserResponse{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "register success")
	resp.Data = createdUser
	return
}

func (u *AuthUsecase) Login(props *model.LoginRequest) (resp model.Response) {
	db := config.GetDatabase()

	existingUser, err := u.repository.GetUserByUsername(props.Username, db)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "invalid username or password")
		return
	} else if err != nil {
		u.log.Errorf("repository.GetUserByUsername: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	if !libs.CheckPasswordHash(props.Password, existingUser.Password) {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := u.jwt.GenerateJWTToken(&existingUser)
	if err != nil {
		u.log.Errorf("libs.GenerateJWTToken: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	data := map[string]any{
		"token": token,
	}

	return model.Response{
		Data:   data,
		Status: libs.CustomResponse(http.StatusOK, "login success"),
	}
}
