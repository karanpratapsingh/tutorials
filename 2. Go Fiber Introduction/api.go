package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type User struct {
	Name    string `json:"name"`
	Twitter string `json:"twitter"`
}

func getUserHander(ctx *fiber.Ctx) error {
	user := User{
		Name:    "Karan",
		Twitter: "karan_6864",
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func createUserHandler(ctx *fiber.Ctx) error {
	body := new(User)
	err := ctx.BodyParser(body)

	if err != nil {
		fmt.Println(err)
		ctx.Status(fiber.StatusBadRequest)
		return err
	}

	user := User{
		Name:    body.Name,
		Twitter: body.Twitter,
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	userApi := app.Group("/user")

	userApi.Get("/", getUserHander)

	userApi.Post("/create", createUserHandler)

	app.Listen(":4000")
}
