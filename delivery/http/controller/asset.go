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

type IController interface {
	GetAssets(ctx *fiber.Ctx) error
}

type Controller struct {
	usecase usecase.IUsecase
	logger  *logrus.Logger
}

func NewController(usecase usecase.IUsecase) IController {
	logger := config.GetLogger()
	return &Controller{
		usecase,
		logger,
	}
}

func (c *Controller) GetAssets(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.AssetRequest
	)

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.logger.Errorf("error parsing request body: %s", err.Error())
		response.Status = libs.CustomResponse(http.StatusBadRequest, "error parsing request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	// validate reqBody struct
	validationErr := libs.ValidateRequest(&reqBody)
	// if there is an error
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		response.Status =
			libs.CustomResponse(http.StatusUnprocessableEntity, "validation error")
		response.Data = errResponse

		return ctx.Status(response.Status.Code).JSON(response)
	}

	response = c.usecase.GetAssets(reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}
