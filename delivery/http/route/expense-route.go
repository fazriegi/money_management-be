package route

import (
	"github.com/fazriegi/money_management-be/delivery/http/controller"
	"github.com/fazriegi/money_management-be/delivery/http/middleware"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/fazriegi/money_management-be/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewExpenseRoute(app *fiber.App, jwt *libs.JWT) {
	repo := repository.NewExpenseRepository()
	usecase := usecase.NewExpenseUsecase(repo)
	controller := controller.NewExpenseController(usecase)

	asset := app.Group("/expenses")
	asset.Get("/", middleware.Authentication(jwt), controller.GetList)
	asset.Put("/", middleware.Authentication(jwt), controller.Update)
}
