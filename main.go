package main

import (
	"github.com/RizkiRdm/go-blog/api"
	"github.com/RizkiRdm/go-blog/db"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// run connection database
	db.Connection()

	// initalize fiber route
	app := fiber.New()

	api.Routes(app)

	app.Listen(":8000")
}
