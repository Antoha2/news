package transport

import (
	"log"

	"github.com/gofiber/fiber"
)

func (a *apiImpl) StartHTTP() error {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
