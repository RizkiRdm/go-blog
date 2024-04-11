package category

import (
	"net/http"

	"github.com/RizkiRdm/go-blog/db"
	"github.com/RizkiRdm/go-blog/pkg/models/category"
	"github.com/gofiber/fiber/v2"
)

func GetCategories(c *fiber.Ctx) error {
	db := db.Connection()
	defer db.Close()
	readCategories := "SELECT `title` FROM `categories`"
	rows, err := db.Query(readCategories)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed get category names",
		})
	}
	defer rows.Close()

	var categories []category.CategoryResponse
	for rows.Next() {
		var category category.CategoryResponse
		if err := rows.Scan(&category.Name); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		categories = append(categories, category)
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": categories,
	})
}
