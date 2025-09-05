package expense

import (
	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *libs.JWT) {
	log := config.GetLogger()
	repo := NewRepository()
	usecase := NewUsecase(log, repo)
	controller := NewController(log, usecase)

	route := app.Group("/expense")
	route.Post("/", middleware.Authentication(jwt), controller.Add)
	route.Get("/category", middleware.Authentication(jwt), controller.ListCategory)
}
