package module

import (
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/auth"
	"github.com/fazriegi/money_management-be/module/expense"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *libs.JWT) {
	auth.NewRoute(app, jwt)
	expense.NewRoute(app, jwt)
}
