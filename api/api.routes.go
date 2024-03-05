package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App) {
	api := r.Group("api")
	V1 := api.Group("V1")

	V1.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("hello")
		return nil
	})
}
