package blog

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/RizkiRdm/go-blog/db"
	"github.com/RizkiRdm/go-blog/pkg/models/blog"
	"github.com/RizkiRdm/go-blog/pkg/models/category"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

// read blogs - GET
func GetBlogs(c *fiber.Ctx) error {
	db := db.Connection()

	// query read blogs
	q := `SELECT 
			blogs.id, 
			users.username, 
			categories.title AS kategori, 
			blogs.thumbnail, 
			blogs.title, 
			blogs.body, 
			blogs.created_at, 
			blogs.updated_at 
			FROM blogs 
			LEFT JOIN users ON blogs.user_id = users.id 
			LEFT JOIN blog_categories ON blogs.id = blog_categories.id_blog
			LEFT JOIN categories ON blog_categories.id_category = categories.id
			GROUP BY blogs.id;
	`

	rows, err := db.Query(q)

	//check err query
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// close connection from db after all data success catching
	defer rows.Close()

	// define blog variable
	var blogs []blog.BlogResponse

	for rows.Next() {
		var blog blog.BlogResponse
		if err := rows.Scan(&blog.Id, &blog.Username, &blog.Category, &blog.Title, &blog.Thumbnail, &blog.Body, &blog.CreatedAt, &blog.UpdatedAt); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		blogs = append(blogs, blog)
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": blogs,
	})
}

// read details blogs - GET
func GetDetailBlog(c *fiber.Ctx) error {
	db := db.Connection()
	id := c.Params("id")
	defer db.Close()

	var blog blog.BlogResponse
	q := `SELECT 
			blogs.id, 
			users.username,
			categories.title AS kategori,
			blogs.thumbnail,
			blogs.title, 
			blogs.body, 
			blogs.created_at,
			blogs.updated_at 
		FROM blogs
			LEFT JOIN users ON blogs.user_id = users.id 
			LEFT JOIN blog_categories ON blogs.id = blog_categories.id_blog 
			LEFT JOIN categories ON blog_categories.id_category = categories.id
		WHERE blogs.id = ?`
	err := db.QueryRow(q, id).Scan(&blog.Id, &blog.Username, &blog.Category, &blog.Title, &blog.Thumbnail, &blog.Body, &blog.CreatedAt, &blog.UpdatedAt)

	if err != nil {
		// jika data kosong
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "blog not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": blog,
	})
}

// create new category - POST
func CreateNewCategory(c *fiber.Ctx) error {
	// handle body parser
	request := new(category.RequestCreateCategory)
	if err := c.BodyParser(request); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "cannot parse",
			"messageErr": err.Error(),
		})
	}

	// define db variable
	db := db.Connection()
	defer db.Close()

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"Message": err.Error(),
		})
	}
	categoryQuery := "INSERT INTO `categories`(`title`) VALUES ('?')"
	result, err := tx.Exec(categoryQuery, request.Name)

	if result != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed insert new category",
			"messageErr": err.Error(),
		})
	}
	// commit transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed commit transaction",
			"messageErr": err.Error(),
		})
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success create new cateogory",
		"data":    request,
	})
}

// create new blog - POST
func CreateBlog(c *fiber.Ctx) error {
	// parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "cannot parse form",
		})
	}

	// convert user_id to int
	userIdStr := form.Value["user_id"][0]
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user id",
		})
	}

	// handle image upload
	file, err := c.FormFile("thumbnail")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message":    "thumbnail upload error",
			"messageErr": err.Error(), // just for debugging
		})
	}

	imageName := uuid.New().String() + filepath.Ext(file.Filename)
	thumbnailPath := filepath.Join("./uploads", imageName)
	if err := c.SaveFile(file, thumbnailPath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed save image",
			"messageErr": err.Error(),
		})
	}

	// define db variable
	db := db.Connection()
	defer db.Close()

	// start transaction
	tx, err := db.Begin()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// insert blog
	blogQuery := "INSERT INTO blogs (user_id, title, thumbnail, body) VALUES (?, ?, ?, ?)"
	result, err := tx.Exec(blogQuery, userId, form.Value["title"][0], thumbnailPath, form.Value["body"][0])
	if err != nil {
		tx.Rollback()

		if err := os.Remove(thumbnailPath); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message":    "failed to rollback image",
				"messageErr": err.Error(), // just debugging
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed insert blog",
			"messageErr": err.Error(),
		})
	}

	blogId, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to retrive blog_id",
		})
	}

	// insert blogId to table blog_categories
	categoryId, err := strconv.Atoi(form.Value["category_id"][0])
	if err != nil {
		tx.Rollback()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message":    "invalid category id",
			"messageErr": err.Error(), // just for debugging
		})
	}
	categoryQuery := "INSERT INTO blog_categories (id_blog, id_category) VALUES (?, ?) ON DUPLICATE KEY UPDATE id_blog = id_blog"
	if _, err := tx.Exec(categoryQuery, blogId, categoryId); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed to insert category",
			"messageErr": err.Error(), // just for debugging
		})
	}

	// commit query
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed commit transaction",
			"messageErr": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    form.Value,
	})
}

// update blog - PATCH
