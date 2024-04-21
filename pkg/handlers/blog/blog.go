package blog

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/RizkiRdm/go-blog/db"
	"github.com/RizkiRdm/go-blog/pkg/models/blog"
	"github.com/RizkiRdm/go-blog/utils"
	"github.com/RizkiRdm/go-blog/utils/dbutil"

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
	slug := c.Params("slug")
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
		WHERE blogs.slug = ?`
	err := db.QueryRow(q, slug).Scan(&blog.Id, &blog.Username, &blog.Category, &blog.Title, &blog.Thumbnail, &blog.Body, &blog.CreatedAt, &blog.UpdatedAt)

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

// create new blog - POST
func CreateBlog(c *fiber.Ctx) error {
	// read jwt token from cookie
	cookies := c.Cookies("jwt")
	if len(cookies) == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "user not authenticated",
		})
	}

	// extract user_id from token
	userId, err := utils.ExtractUserId(cookies)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message":    "invalid token for user_id",
			"messageErr": err.Error(),
		})
	}

	// parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "cannot parse form",
		})
	}

	// handle upload image
	thumbnailPath, err := dbutil.UploadImage(c, "thumbnail")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed to upload image",
			"messageErr": err.Error(),
		})
	}

	// start transaction
	db := db.Connection()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return dbutil.HandleTransactionError(tx, c, err)
	}
	defer tx.Rollback()

	// insert blog
	if err := dbutil.CreateBlog(tx, form, userId, thumbnailPath); err != nil {
		return dbutil.HandleTransactionError(tx, c, err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return dbutil.HandleTransactionError(tx, c, err)
	}

	// return data
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    form.Value,
	})
}

// update blog - PATCH
func UpdateBlog(c *fiber.Ctx) error {
	// read jwt token from cookie
	cookies := c.Cookies("jwt")
	if len(cookies) == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "user not authenticated",
		})
	}

	// extract user_id from token
	userId, err := utils.ExtractUserId(cookies)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message":    "invalid token for user_id",
			"messageErr": err.Error(),
		})
	}

	// parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "cannot parse form",
		})
	}

	// handle image upload
	thumbnailPath, err := dbutil.UploadImage(c, "thumbnail")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "error save image",
			"messageErr": err.Error(),
		})
	}

	// start database transaction
	db := db.Connection()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return dbutil.HandleTransactionError(tx, c, err)
	}
	defer tx.Rollback()

	// delete old thumbnail
	id := c.Params("id")
	oldThumbnail, err := dbutil.GetOldThumbnail(tx, id, userId)
	if err != nil {
		return dbutil.HandleTransactionError(tx, c, err)
	}

	if err := os.Remove(oldThumbnail); err != nil {
		return dbutil.HandleTransactionError(tx, c, err)
	}

	// update blog
	if err := dbutil.UpdateBlog(tx, form, id, userId, thumbnailPath); err != nil {
		return dbutil.HandleTransactionError(tx, c, err)
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return dbutil.HandleTransactionError(tx, c, err)
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success update",
		"data":    form.Value,
	})
}

// SEARCH BLOG - GET
func SearchBlog(c *fiber.Ctx) error {
	searchQuery := c.Query("q")

	// check query not empty
	if len(searchQuery) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "search query empty",
		})
	}

	db := db.Connection()
	defer db.Close()

	search := "SELECT id, title, thumbnail, body FROM `blogs` WHERE `title` LIKE ? OR `body` LIKE ?"
	rows, err := db.Query(search, "%"+searchQuery+"%", "%"+searchQuery+"%")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed to search blog",
			"messageErr": err.Error(),
		})
	}
	defer rows.Close()

	var searchResults []map[string]interface{}

	for rows.Next() {
		var (
			id        int
			title     string
			thumbnail string
			body      string
		)
		if err := rows.Scan(&id, &title, &thumbnail, &body); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message":    "failed to scan row",
				"messageErr": err.Error(),
			})
		}

		blog := map[string]interface{}{
			"id":        id,
			"title":     title,
			"thumbnail": thumbnail,
			"body":      body,
		}
		searchResults = append(searchResults, blog)
	}

	// if blog not found
	if len(searchResults) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": fmt.Sprintf("no blog with query : '%s'", searchQuery),
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": searchResults,
	})
}
