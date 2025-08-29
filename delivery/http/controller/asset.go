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

type IAssetController interface {
	GetAssets(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type AssetController struct {
	usecase usecase.IAssetUsecase
	logger  *logrus.Logger
}

func NewAssetController(usecase usecase.IAssetUsecase) IAssetController {
	logger := config.GetLogger()
	return &AssetController{
		usecase,
		logger,
	}
}

func (c *AssetController) GetAssets(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.AssetRequest
		user     = ctx.Locals("user").(model.User)
	)

	if err := ctx.QueryParser(&reqBody); err != nil {
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

	response = c.usecase.GetAssets(&user, &reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *AssetController) Update(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.InsertAssetRequest
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
