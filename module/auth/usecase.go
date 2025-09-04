package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	globalModel "github.com/fazriegi/money_management-be/model"
	"github.com/fazriegi/money_management-be/module/auth/model"
	"github.com/fazriegi/money_management-be/module/master/user"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"golang.org/x/sync/errgroup"

	"github.com/sirupsen/logrus"
)

type Usecase interface {
	Register(props *model.RegisterRequest) (resp globalModel.Response)
	Login(props *model.LoginRequest) (resp globalModel.Response)
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

func (u *usecase) Register(props *model.RegisterRequest) (resp globalModel.Response) {
	var (
		err            error
		user           userModel.User
		hashedPassword string
		db             = config.GetDatabase()
	)

	user, err = u.repository.GetByUsername(props.Username, db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		u.log.Errorf("repository.GetByUsername: %s", err.Error())
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

	tx, err := db.Beginx()
	if err != nil {
		u.log.Errorf("error start transaction: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
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
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	g, _ := errgroup.WithContext(context.Background())
	if err != nil {
		u.log.Errorf("failed start goroutine: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
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
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	if err := tx.Commit(); err != nil {
		u.log.Errorf("failed commit tx: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	createdUser := userModel.UserResponse{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
	}

	resp.Status = libs.CustomResponse(http.StatusCreated, "register success")
	resp.Data = createdUser
	return
}

func (u *usecase) Login(props *model.LoginRequest) (resp globalModel.Response) {
	db := config.GetDatabase()

	existingUser, err := u.repository.GetByUsername(props.Username, db)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "invalid username or password")
		return
	} else if err != nil {
		u.log.Errorf("repository.GetByUsername: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	if !libs.CheckPasswordHash(props.Password, existingUser.Password) {
		resp.Status = libs.CustomResponse(http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := u.jwt.GenerateJWTToken(existingUser.ID, existingUser.Email, existingUser.Username)
	if err != nil {
		u.log.Errorf("libs.GenerateJWTToken: %s", err.Error())
		resp.Status = libs.CustomResponse(http.StatusInternalServerError, constant.ServerErr)
		return
	}

	data := map[string]any{
		"token": token,
		"user": userModel.UserResponse{
			Name:     existingUser.Name,
			Username: existingUser.Username,
		},
	}

	return globalModel.Response{
		Data:   data,
		Status: libs.CustomResponse(http.StatusOK, "login success"),
	}
}
