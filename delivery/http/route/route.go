package route

import (
	"github.com/fazriegi/money_management-be/libs"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *libs.JWT) {
	NewAssetRoute(app)
	NewIncomeRoute(app)
	NewExpenseRoute(app)
	NewLiabilityRoute(app)
	NewAuthRoute(app, jwt)
}
