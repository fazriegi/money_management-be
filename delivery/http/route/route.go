package route

import (
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App) {
	NewUserRoute(app)
}
