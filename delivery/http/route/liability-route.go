package route

import (
	"github.com/fazriegi/money_management-be/delivery/http/controller"
	"github.com/fazriegi/money_management-be/delivery/http/middleware"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/fazriegi/money_management-be/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewLiabilityRoute(app *fiber.App, jwt *libs.JWT) {
	repo := repository.NewLiabilityRepository()
	usecase := usecase.NewLiabilityUsecase(repo)
	controller := controller.NewLiabilityController(usecase)

	asset := app.Group("/liabilities")
	asset.Get("/", middleware.Authentication(jwt), controller.GetList)
	asset.Put("/", middleware.Authentication(jwt), controller.Update)
}
