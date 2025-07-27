package route

import (
	"github.com/fazriegi/money_management-be/delivery/http/controller"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/fazriegi/money_management-be/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAssetRoute(app *fiber.App) {
	repo := repository.NewRepository()
	usecase := usecase.NewUsecase(repo)
	controller := controller.NewController(usecase)

	asset := app.Group("/assets")
	asset.Get("/", controller.GetAssets)
}
