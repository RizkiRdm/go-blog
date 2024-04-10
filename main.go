package main

import (
	"github.com/RizkiRdm/go-blog/api"
	"github.com/RizkiRdm/go-blog/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// run connection database
	db.Connection()

	// initalize fiber route
	app := fiber.New()

	// setup CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowOrigins:     "*",
	}))

	api.Routes(app)

	// app.Use(middleware.Midleware())
	app.Listen(":8000")
}
