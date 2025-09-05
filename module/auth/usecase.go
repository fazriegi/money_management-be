package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/auth/model"
	"github.com/fazriegi/money_management-be/module/common"
	"github.com/fazriegi/money_management-be/module/master/user"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"golang.org/x/sync/errgroup"

	"github.com/sirupsen/logrus"
)

type Usecase interface {
	Register(props *model.RegisterRequest) (resp common.Response)
	Login(props *model.LoginRequest) (resp common.Response)
}

type usecase struct {
	repository user.Repository
	log        *logrus.Logger
	jwt        *libs.JWT
}

func NewUsecase(repository user.Repository, jwt *libs.JWT) Usecase {
	log := config.GetLogger()

	return &usecase{
		repository,
		log,
		jwt,
	}
}

func (u *usecase) Register(props *model.RegisterRequest) (resp common.Response) {
	var (
		err            error
		user           userModel.User
		hashedPassword string
		db             = config.GetDatabase()
	)

	user, err = u.repository.GetByUsername(props.Username, db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.log.Errorf("repository.GetByUsername: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	if user.Username != "" {
		return resp.CustomResponse(http.StatusBadRequest, "username already exists", nil)
	}

	if hashedPassword, err = libs.HashPassword(props.Password); err != nil {
		u.log.Errorf("libs.HashPassword: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error start transaction: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}
	defer tx.Rollback()

	user = userModel.User{
		Name:     props.Name,
		Email:    props.Email,
		Password: hashedPassword,
		Username: props.Username,
	}

	userId, err := u.repository.Insert(&user, tx)
	if err != nil {
		u.log.Errorf("repository.Insert: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	g, _ := errgroup.WithContext(context.Background())
	if err != nil {
		u.log.Errorf("failed start goroutine: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	g.Go(func() error {
		err = u.repository.CreateExpenseCat(userId, tx)
		if err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		err = u.repository.CreateIncomeCat(userId, tx)
		if err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		err = u.repository.CreateAssetCat(userId, tx)
		if err != nil {
			return err
		}

		return nil
	})

	err = g.Wait()
	if err != nil {
		u.log.Errorf("failed initialized category: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	if err := tx.Commit(); err != nil {
		u.log.Errorf("failed commit tx: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	createdUser := userModel.UserResponse{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
	}

	return resp.CustomResponse(http.StatusCreated, "success", createdUser)
}

func (u *usecase) Login(props *model.LoginRequest) (resp common.Response) {
	db := config.GetDatabase()

	existingUser, err := u.repository.GetByUsername(props.Username, db)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return resp.CustomResponse(http.StatusUnauthorized, "invalid username or password", nil)
	} else if err != nil {
		u.log.Errorf("repository.GetByUsername: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	if !libs.CheckPasswordHash(props.Password, existingUser.Password) {
		return resp.CustomResponse(http.StatusUnauthorized, "invalid username or password", nil)
	}

	token, err := u.jwt.GenerateJWTToken(existingUser.ID, existingUser.Email, existingUser.Username)
	if err != nil {
		u.log.Errorf("libs.GenerateJWTToken: %s", err.Error())
		return resp.CustomResponse(http.StatusInternalServerError, constant.ServerErr, nil)
	}

	data := map[string]any{
		"token": token,
		"user": userModel.UserResponse{
			Name:     existingUser.Name,
			Username: existingUser.Username,
		},
	}

	return resp.CustomResponse(http.StatusOK, "success", data)
}
