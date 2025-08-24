package controller

import (
	"net/http"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/model"
	"github.com/fazriegi/money_management-be/usecase"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

type IAuthController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type AuthController struct {
	usecase usecase.IAuthUsecase
	logger  *logrus.Logger
}

func NewAuthController(usecase usecase.IAuthUsecase) IAuthController {
	logger := config.GetLogger()
	return &AuthController{
		usecase,
		logger,
	}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.RegisterRequest
	)

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.logger.Errorf("error parsing request body: %s", err.Error())
		response.Status = libs.CustomResponse(http.StatusBadRequest, "error parsing request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	// validate reqBody struct
	validationErr := libs.ValidateRequest(&reqBody)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "validation error")
		response.Data = errResponse

		return ctx.Status(response.Status.Code).JSON(response)
	}

	response = c.usecase.Register(&reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.LoginRequest
	)

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.logger.Errorf("error parsing request body: %s", err.Error())
		response.Status = libs.CustomResponse(http.StatusBadRequest, "error parsing request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	// validate reqBody struct
	validationErr := libs.ValidateRequest(&reqBody)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "validation error")
		response.Data = errResponse

		return ctx.Status(response.Status.Code).JSON(response)
	}

	response = c.usecase.Login(&reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}
