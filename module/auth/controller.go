package auth

import (
	"net/http"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	entityModel "github.com/fazriegi/money_management-be/model"
	"github.com/fazriegi/money_management-be/module/auth/model"
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

func NewWController(usecase Usecase) Controller {
	logger := config.GetLogger()
	return &controller{
		usecase,
		logger,
	}
}

func (c *controller) Register(ctx *fiber.Ctx) error {
	var (
		response entityModel.Response
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

func (c *controller) Login(ctx *fiber.Ctx) error {
	var (
		response entityModel.Response
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
