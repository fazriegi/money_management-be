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

type IExpenseController interface {
	GetList(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type ExpenseController struct {
	usecase usecase.IExpenseUsecase
	logger  *logrus.Logger
}

func NewExpenseController(usecase usecase.IExpenseUsecase) IExpenseController {
	logger := config.GetLogger()
	return &ExpenseController{
		usecase,
		logger,
	}
}

func (c *ExpenseController) GetList(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.GetExpenseRequest
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

func (c *ExpenseController) Update(ctx *fiber.Ctx) error {
	var (
		response model.Response
		reqBody  model.UpdateExpenseRequest
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
