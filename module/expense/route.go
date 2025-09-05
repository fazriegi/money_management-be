package expense

import (
	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/delivery/http/middleware"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *libs.JWT) {
	log := config.GetLogger()
	repo := NewRepository()
	usecase := NewUsecase(log, repo)
	controller := NewController(log, usecase)

	route := app.Group("/expense")
	route.Post("/", middleware.Authentication(jwt), controller.Add)
}
