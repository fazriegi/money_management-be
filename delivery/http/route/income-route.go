package route

import (
	"github.com/fazriegi/money_management-be/delivery/http/controller"
	"github.com/fazriegi/money_management-be/delivery/http/middleware"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/fazriegi/money_management-be/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewIncomeRoute(app *fiber.App, jwt *libs.JWT) {
	repo := repository.NewIncomeRepository()
	usecase := usecase.NewIncomeUsecase(repo)
	controller := controller.NewIncomeController(usecase)

	asset := app.Group("/incomes")
	asset.Get("/", middleware.Authentication(jwt), controller.GetList)
	asset.Put("/", middleware.Authentication(jwt), controller.Update)
}
