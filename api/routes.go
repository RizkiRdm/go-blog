package api

import (
	"github.com/RizkiRdm/go-blog/pkg/handlers/blog"
	"github.com/RizkiRdm/go-blog/pkg/handlers/users"
	"github.com/gofiber/fiber/v2"
)

func Routes(r *fiber.App) {
	api := r.Group("api")
	V1 := api.Group("v1")

	// AUTH ROUTES
	// LOGIN USER
	V1.Post("/login", users.LoginUser)
	// REGISTER USER
	V1.Post("/register", users.RegisterUser)
	// LOGOUT USER
	V1.Post("/logout", users.LogoutUser)

	// BLOG ROUTES
	// GET ALL BLOGS
	V1.Get("/blogs", blog.GetBlogs)
	// GET DETAIL BLOG
	V1.Get("/blogs/:id", blog.GetDetailBlog)
	// CREATE NEW BLOG
	V1.Post("/blogs", blog.CreateBlog)
	// UPDATE BLOG
	V1.Patch("/blogs/:id", blog.UpdateBlog)
	// CREATE NEW CATEGORY
	V1.Post("/categories", blog.CreateNewCategory)
}
