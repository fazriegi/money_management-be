package main

import (
	"fmt"
	"log"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/libs"
	"github.com/fazriegi/money_management-be/middleware"
	"github.com/fazriegi/money_management-be/module"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	viperConfig := config.NewViper()
	config.NewDatabase(viperConfig)
	jwt := libs.InitJWT(viperConfig)
	file := config.NewLogger(viperConfig)
	defer file.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(middleware.LogMiddleware())
	port := viperConfig.GetInt("web.port")
	module.NewRoute(app, jwt)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
