package route

import (
	"github.com/fazriegi/money_management-be/delivery/http/controller"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/repository"
	"github.com/fazriegi/money_management-be/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAuthRoute(app *fiber.App, jwt *libs.JWT) {
	repo := repository.NewUserRepository()
	usecase := usecase.NewAuthUsecase(repo, jwt)
	controller := controller.NewAuthController(usecase)

	Auth := app.Group("/auth")
	Auth.Post("/register", controller.Register)
	Auth.Post("/login", controller.Login)
}
