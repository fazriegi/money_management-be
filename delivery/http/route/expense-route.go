package route

import (
	"github.com/fazriegi/money_management-be/delivery/http/controller"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/fazriegi/money_management-be/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewExpenseRoute(app *fiber.App) {
	repo := repository.NewExpenseRepository()
	usecase := usecase.NewExpenseUsecase(repo)
	controller := controller.NewExpenseController(usecase)

	asset := app.Group("/expenses")
	asset.Get("/", controller.GetList)
	asset.Put("/", controller.Update)
}
