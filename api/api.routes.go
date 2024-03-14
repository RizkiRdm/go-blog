package api

import (
	"github.com/RizkiRdm/go-blog/pkg/handlers"
	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App) {
	api := r.Group("api")
	V1 := api.Group("v1")

	// GET ALL BLOGS
	V1.Get("/blogs", handlers.GetBlogs)
	// GET DETAIL BLOG
	V1.Get("/blogs/:id", handlers.GetDetailBlog)
	// CREATE NEW BLOG
	V1.Post("/blogs", handlers.CreateBlog)
	// UPDATE CURRENT BLOG
	V1.Post("/blogs/:id/:userId", handlers.UpdateBlog)
}
