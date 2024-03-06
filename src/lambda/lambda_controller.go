package lambda

import (
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"src/server"
)

func SetController() {
	server.App.Get("/lambda", func(c fiber.Ctx) error {
		var lambda Lambda
		var id string
		var primitiveId primitive.ObjectID
		var err error
		if id = c.Query("id", ""); id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid parameter id"})
		}
		if primitiveId, err = primitive.ObjectIDFromHex(id); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if lambda, err = getLambda(primitiveId); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(lambda)
	})

	server.App.Post("/lambda", func(c fiber.Ctx) error {
		var lambda Lambda
		var id primitive.ObjectID
		var err error
		if err = json.Unmarshal(c.Body(), &lambda); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err = validateInsertLambda(&lambda); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if id, err = insertLambda(&lambda); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(id.Hex())
	})
}
