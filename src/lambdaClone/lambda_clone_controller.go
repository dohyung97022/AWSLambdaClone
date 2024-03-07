package lambdaClone

import (
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"src/server"
)

func SetController() {
	server.App.Get("/lambda", func(c fiber.Ctx) error {
		var err error
		var primitiveId *primitive.ObjectID
		if primitiveId, err = validateGetLambda(&c); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		var lambda *Lambda
		if lambda, err = getLambda(primitiveId); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		var code string
		if code, err = s3Client.readCodeFromS3(lambda); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		lambda.Code = code
		return c.Status(fiber.StatusOK).JSON(lambda)
	})

	server.App.Get("/lambda/list", func(c fiber.Ctx) error {
		var err error
		var option *options.FindOptions
		if option, err = validateGetLambdas(&c); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		var lambdas *[]Lambda
		if lambdas, err = getLambdas(option); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		var cnt int64
		var pages int64
		if cnt, pages, err = getLambdasCnt(option.Limit); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"lambdas": lambdas, "pages": pages, "cnt": cnt})
	})

	server.App.Post("/lambda", func(c fiber.Ctx) error {
		var err error
		var lambda *Lambda
		if lambda, err = validateInsertLambda(&c); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		var id *primitive.ObjectID
		if id, err = insertLambda(lambda); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		lambda.Id = id.Hex()
		if err = s3Client.uploadCodeToS3(lambda); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(id.Hex())
	})

	server.App.Patch("/lambda", func(c fiber.Ctx) error {
		var err error
		var lambda *Lambda
		if lambda, err = validateUpdateLambda(&c); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err = updateLambda(lambda); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if err = s3Client.uploadCodeToS3(lambda); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})

	server.App.Delete("/lambda", func(c fiber.Ctx) error {
		var err error
		var id *primitive.ObjectID
		if id, err = validateDeleteLambda(&c); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err = deleteLambda(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})
}
