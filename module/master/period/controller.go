package period

import (
	"github.com/fazriegi/money_management-be/module/common"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	Get(ctx *fiber.Ctx) error
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

func (c *controller) Get(ctx *fiber.Ctx) error {
	var (
		response common.Response
		user     = ctx.Locals("user").(userModel.User)
	)

	response = c.usecase.Get(&user)

	return ctx.Status(response.Status.Code).JSON(response)
}
