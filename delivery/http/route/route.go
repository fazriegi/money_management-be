package route

import (
	"github.com/fazriegi/money_management-be/libs"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *libs.JWT) {
	NewAssetRoute(app, jwt)
	NewIncomeRoute(app, jwt)
	NewExpenseRoute(app, jwt)
	NewLiabilityRoute(app, jwt)
	NewAuthRoute(app, jwt)
}
