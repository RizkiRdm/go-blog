package api

import (
	blog "github.com/RizkiRdm/go-blog/pkg/handlers"
	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App) {
	api := r.Group("api")
	V1 := api.Group("v1")

	// GET ALL BLOGS
	V1.Get("/blogs", blog.GetBlogs)
	// GET DETAIL BLOG
	V1.Get("/blogs/:id", blog.GetDetailBlog)
	// CREATE NEW BLOG
	// V1.Post("/blogs", middleware.AuthMiddleware(), blog.HandleBlogWithDetails)
	V1.Post("/blogs", blog.CreateBlog)
}
