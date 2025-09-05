package auth

import (
	"net/http"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/auth/model"
	"github.com/fazriegi/money_management-be/module/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type controller struct {
	usecase Usecase
	logger  *logrus.Logger
}

func NewController(usecase Usecase) Controller {
	logger := config.GetLogger()
	return &controller{
		usecase,
		logger,
	}
}

func (c *controller) Register(ctx *fiber.Ctx) error {
	var (
		response common.Response
		reqBody  model.RegisterRequest
	)

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.logger.Errorf("error parsing request body: %s", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(response.CustomResponse(http.StatusBadRequest, constant.ParseReqBodyErr, nil))
	}

	// validate reqBody struct
	validationErr := libs.ValidateRequest(&reqBody)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		return ctx.Status(response.Status.Code).JSON(response.CustomResponse(http.StatusUnprocessableEntity, constant.ValidationErr, errResponse))
	}

	response = c.usecase.Register(&reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *controller) Login(ctx *fiber.Ctx) error {
	var (
		response common.Response
		reqBody  model.LoginRequest
	)

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.logger.Errorf("error parsing request body: %s", err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(response.CustomResponse(http.StatusBadRequest, constant.ParseReqBodyErr, nil))
	}

	// validate reqBody struct
	validationErr := libs.ValidateRequest(&reqBody)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		return ctx.Status(response.Status.Code).JSON(response.CustomResponse(http.StatusUnprocessableEntity, constant.ValidationErr, errResponse))
	}

	response = c.usecase.Login(&reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}
