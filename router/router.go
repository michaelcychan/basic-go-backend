package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michaelcychan/basic-go-backend/handler"
)

func SetupRouter(app *fiber.App) {
	homePath := app.Group("/home")
	kingdomPath := app.Group("/kingdom")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	homePath.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome Home!")
	})

	homePath.Get("/get-temp/:roomname", handler.IndexHandler)

	kingdomPath.Get("/", handler.KingdomRootHandler)

	kingdomPath.Get("/get-all-monarch", handler.FindAllMonarchHandler)

	kingdomPath.Get("/get-monarch/:monarch", handler.FindMonarchHandler)
}
