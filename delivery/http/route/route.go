package route

import (
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App) {
	NewAssetRoute(app)
	NewIncomeRoute(app)
	NewExpenseRoute(app)
}
