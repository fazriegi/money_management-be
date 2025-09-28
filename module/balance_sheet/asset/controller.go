package asset

import (
	"net/http"

	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/balance_sheet/asset/model"
	"github.com/fazriegi/money_management-be/module/common"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Add(ctx *fiber.Ctx) error
	ListCategory(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type controller struct {
	log     *logrus.Logger
	usecase Usecase
}

func NewController(log *logrus.Logger, usecase Usecase) Controller {
	return &controller{
		log,
		usecase,
	}
}

func (c *controller) Add(ctx *fiber.Ctx) error {
	var (
		response common.Response
		reqBody  model.AddRequest

		user = ctx.Locals("user").(userModel.User)
	)

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.log.Errorf("error parsing request body: %s", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(response.CustomResponse(http.StatusBadRequest, constant.ParseReqBodyErr, nil))
	}

	// validate reqBody struct
	validationErr := libs.ValidateRequest(&reqBody)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		return ctx.Status(http.StatusUnprocessableEntity).JSON(response.CustomResponse(http.StatusUnprocessableEntity, constant.ValidationErr, errResponse))
	}

	response = c.usecase.Add(&user, &reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *controller) ListCategory(ctx *fiber.Ctx) error {
	var (
		response common.Response
		user     = ctx.Locals("user").(userModel.User)
	)

	response = c.usecase.ListCategory(&user)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *controller) List(ctx *fiber.Ctx) error {
	var (
		response common.Response
		reqBody  model.ListRequest

		user = ctx.Locals("user").(userModel.User)
	)

	if err := ctx.QueryParser(&reqBody); err != nil {
		c.log.Errorf("error parsing query param: %s", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(response.CustomResponse(http.StatusBadRequest, constant.ParseQueryParamErr, nil))
	}

	response = c.usecase.List(&user, &reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *controller) Update(ctx *fiber.Ctx) error {
	var (
		response common.Response
		reqBody  model.UpdateRequest

		user = ctx.Locals("user").(userModel.User)
	)

	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.log.Errorf("error get param: %s", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(response.CustomResponse(http.StatusBadRequest, "invalid id", nil))
	}

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.log.Errorf("error parsing request body: %s", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(response.CustomResponse(http.StatusBadRequest, constant.ParseReqBodyErr, nil))
	}

	// validate reqBody struct
	validationErr := libs.ValidateRequest(&reqBody)
	if len(validationErr) > 0 {
		errResponse := map[string]any{
			"errors": validationErr,
		}

		return ctx.Status(http.StatusUnprocessableEntity).JSON(response.CustomResponse(http.StatusUnprocessableEntity, constant.ValidationErr, errResponse))
	}

	reqBody.ID = uint(id)
	response = c.usecase.Update(&user, &reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}

func (c *controller) Delete(ctx *fiber.Ctx) error {
	var (
		response common.Response
		user     = ctx.Locals("user").(userModel.User)
	)

	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.log.Errorf("error get param: %s", err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(response.CustomResponse(http.StatusBadRequest, "invalid id", nil))
	}

	response = c.usecase.Delete(&user, uint(id))

	return ctx.Status(response.Status.Code).JSON(response)
}
