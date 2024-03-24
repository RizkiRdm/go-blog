package handlers

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RizkiRdm/go-blog/db"
	"github.com/RizkiRdm/go-blog/pkg/models/blog"

	"github.com/gofiber/fiber/v2"
)

func SplitTags(tagString string) []string {
	if tagString == "" {
		return nil
	}

	return strings.Split(tagString, ",")
}

// read blogs - GET
func GetBlogs(c *fiber.Ctx) error {
	db := db.Connection()

	// query read blogs
	q := `SELECT 
			blogs.id, 
			users.username, 
			categories.title AS kategori, 
			GROUP_CONCAT(tags.name) AS tag, 
			blogs.thumbnail, 
			blogs.title, 
			blogs.body, 
			blogs.created_at, 
			blogs.updated_at 
			FROM blogs 
			LEFT JOIN users ON blogs.user_id = users.id 
			LEFT JOIN blog_categories ON blogs.id = blog_categories.id_blog
			LEFT JOIN categories ON blog_categories.id_category = categories.id
			LEFT JOIN blog_tags ON blogs.id = blog_tags.id_blog 
			LEFT JOIN tags ON blog_tags.id_tag = tags.id  
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
		var tags string
		if err := rows.Scan(&blog.Id, &blog.Username, &blog.Category, &tags, &blog.Title, &blog.Thumbnail, &blog.Body, &blog.CreatedAt, &blog.UpdatedAt); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		blog.TagName = SplitTags(tags)
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
	var tag string
	q := "SELECT blogs.id, users.username, categories.title AS kategori, GROUP_CONCAT(DISTINCT tags.name) AS tag, blogs.thumbnail, blogs.title, blogs.body, blogs.created_at, blogs.updated_at FROM blogs LEFT JOIN users ON blogs.user_id = users.id LEFT JOIN blog_categories ON blogs.id = blog_categories.id_blog LEFT JOIN categories ON blog_categories.id_category = categories.id LEFT join blog_tags ON blogs.id = blog_tags.id_tag LEFT JOIN tags ON blogs.id = tags.id WHERE blogs.id = ?"
	err := db.QueryRow(q, id).Scan(&blog.Id, &blog.Username, &blog.Category, &tag, &blog.Title, &blog.Thumbnail, &blog.Body, &blog.CreatedAt, &blog.UpdatedAt)

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

	blog.TagName = SplitTags(tag)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": blog,
	})
}

// create new blog - POST
func CreateBlog(c *fiber.Ctx) error {
	blogReq := new(blog.RequestCreateBlog)
	if err := c.BodyParser(blogReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// convert user_id from string to int
	userId, err := strconv.Atoi(blogReq.UsernameId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid user ID",
		})
	}

	// handle image upload
	file, err := c.FormFile("thumbnail")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "thumbnail upload error",
		})
	}

	// save image to uploads directory
	thumbnailPath := filepath.Join("./uploads", file.Filename)
	if err := c.SaveFile(file, thumbnailPath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed save image",
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
	blogQuery := "INSERT INTO `blogs` (`user_id`, `title`, `thumbnail`, `body`) VALUES (?, ?, ?, ?)"
	result, err := tx.Exec(blogQuery, userId, blogReq.Title, thumbnailPath, blogReq.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed insert blog",
		})
	}
	blogId, err := result.LastInsertId()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve blog ID",
		})
	}

	// insert blogId to table blog_categories
	// convert category_id from string to int
	categoryId, err := strconv.Atoi(blogReq.CategoryId)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid category ID",
		})
	}
	categoryQuery := "INSERT INTO `blog_categories` (`id_blog`, `id_category`) VALUES (?, ?) ON DUPLICATE KEY UPDATE id_blog = id_blog"
	if _, err := tx.Exec(categoryQuery, blogId, categoryId); err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to insert category",
		})
	}

	// insert blog_id to blog_tags
	tagQuery := "INSERT INTO `blog_tags` (`id_blog`, `id_tag`) VALUES (?, ?) ON DUPLICATE KEY UPDATE id_blog = id_blog"
	for _, tagStr := range blogReq.TagsId {
		// convert tags_id to int
		tagId, err := strconv.Atoi(tagStr)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid tags id",
			})
		}
		if _, err := tx.Exec(tagQuery, blogId, tagId); err != nil {
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to insert tags"})
		}
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    blogReq,
	})
}
