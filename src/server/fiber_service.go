package server

import (
	"github.com/gofiber/fiber/v3"
)

var App = fiber.New()

func Listen() error {
	err := App.Listen(":8081")
	return err
}
