package cashflow

import (
	"net/http"

	"github.com/fazriegi/money_management-be/constant"
	"github.com/fazriegi/money_management-be/module/cashflow/model"
	"github.com/fazriegi/money_management-be/module/common"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	List(ctx *fiber.Ctx) error
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
