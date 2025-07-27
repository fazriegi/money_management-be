package route

import (
	"github.com/fazriegi/money_management-be/delivery/http/controller"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/fazriegi/money_management-be/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewIncomeRoute(app *fiber.App) {
	repo := repository.NewIncomeRepository()
	usecase := usecase.NewIncomeUsecase(repo)
	controller := controller.NewIncomeController(usecase)

	asset := app.Group("/incomes")
	asset.Get("/", controller.GetList)
	asset.Put("/", controller.Update)
}
