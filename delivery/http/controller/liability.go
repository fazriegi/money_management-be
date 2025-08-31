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

type ILiabilityController interface {
	GetList(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	ValidateDelete(ctx *fiber.Ctx) error
}

type LiabilityController struct {
	usecase usecase.ILiabilityUsecase
	logger  *logrus.Logger
}

func NewLiabilityController(usecase usecase.ILiabilityUsecase) ILiabilityController {
	logger := config.GetLogger()
	return &LiabilityController{
		usecase,
		logger,
	}
}

func (c *LiabilityController) GetList(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.GetLiabilityRequest
		user     = ctx.Locals("user").(model.User)
	)

	if err := ctx.QueryParser(&reqBody); err != nil {
		c.logger.Errorf("error parsing request query param: %s", err.Error())
		response.Status = libs.CustomResponse(http.StatusBadRequest, "error parsing request query param")
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

	response = c.usecase.GetList(&user, &reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *LiabilityController) Update(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.UpdateLiabilityRequest
		user     = ctx.Locals("user").(model.User)
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

	response = c.usecase.Update(&user, &reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *LiabilityController) ValidateDelete(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.ValidateDeleteRequest
		user     = ctx.Locals("user").(model.User)
	)

	if err := ctx.QueryParser(&reqBody); err != nil {
		c.logger.Errorf("error parsing request query param: %s", err.Error())
		response.Status = libs.CustomResponse(http.StatusBadRequest, "error parsing request query param")
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

	response = c.usecase.ValidateDelete(&user, &reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}
