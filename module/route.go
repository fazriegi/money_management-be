package module

import (
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/auth"
	"github.com/fazriegi/money_management-be/module/cashflow"
	"github.com/fazriegi/money_management-be/module/master/period"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *libs.JWT) {
	auth.NewRoute(app, jwt)
	cashflow.NewRoute(app, jwt)
	period.NewRoute(app, jwt)
}
