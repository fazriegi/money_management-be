package route

import (
	"github.com/fazriegi/money_management-be/delivery/http/controller"
	"github.com/fazriegi/money_management-be/delivery/http/middleware"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/fazriegi/money_management-be/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAssetRoute(app *fiber.App, jwt *libs.JWT) {
	repo := repository.NewAssetRepository()
	usecase := usecase.NewAssetUsecase(repo)
	controller := controller.NewAssetController(usecase)

	asset := app.Group("/assets")
	asset.Get("/", middleware.Authentication(jwt), controller.GetAssets)
	asset.Put("/", middleware.Authentication(jwt), controller.Update)
}
