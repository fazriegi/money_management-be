package main

import (
	"fmt"
	"log"

	"github.com/fazriegi/money_management-be/config"
	"github.com/fazriegi/money_management-be/delivery/http/middleware"
	"github.com/fazriegi/money_management-be/delivery/http/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	viperConfig := config.NewViper()
	config.NewDatabase(viperConfig)
	file := config.NewLogger(viperConfig)
	defer file.Close()

	app := fiber.New()

	app.Use(middleware.LogMiddleware())
	port := viperConfig.GetInt("web.port")
	route.NewRoute(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
