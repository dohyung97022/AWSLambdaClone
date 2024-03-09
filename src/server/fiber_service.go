package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

var App = fiber.New()

func init() {
	App.Use(cors.New(cors.Config{
		AllowOrigins: "https://www.dev-doe.com, https://dev-doe.com, http://localhost:8080",
	}))
}

func Listen() error {
	err := App.Listen(":443")
	return err
}
