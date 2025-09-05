package expense

import (
	"net/http"

	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/expense/model"
	userModel "github.com/fazriegi/money_management-be/module/master/user/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	globalModel "github.com/fazriegi/money_management-be/model"
)

type Controller interface {
	Add(ctx *fiber.Ctx) error
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
		response globalModel.Response
		reqBody  model.AddRequest

		user = ctx.Locals("user").(userModel.User)
	)

	if err := ctx.BodyParser(&reqBody); err != nil {
		c.log.Errorf("error parsing request body: %s", err.Error())
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

	response = c.usecase.Add(&user, &reqBody)

	return ctx.Status(response.Status.Code).JSON(response)
}
