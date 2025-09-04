package auth

import (
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/module/master/user"

	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, jwt *libs.JWT) {
	repo := user.NewRepository()
	usecase := NewUsecase(repo, jwt)
	controller := NewWController(usecase)

	Auth := app.Group("/auth")
	Auth.Post("/register", controller.Register)
	Auth.Post("/login", controller.Login)
}
