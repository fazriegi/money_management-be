package cashflow

import (
	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/middleware"
	"github.com/fazriegi/money_management-be/module/cashflow/expense"
	"github.com/fazriegi/money_management-be/module/cashflow/income"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *libs.JWT) {
	expense.NewRoute(app, jwt)
	income.NewRoute(app, jwt)
	cashflowRoute(app, jwt)
}

func cashflowRoute(app *fiber.App, jwt *libs.JWT) {
	log := config.GetLogger()
	expenseRepo := expense.NewRepository()
	incomeRepo := income.NewRepository()

	repo := NewRepository(expenseRepo, incomeRepo)
	usecase := NewUsecase(log, repo)
	controller := NewController(log, usecase)

	route := app.Group("/cashflow")
	route.Get("/", middleware.Authentication(jwt), controller.List)
}
