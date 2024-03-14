package handlers

import (
	"database/sql"
	"net/http"
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
	GROUP_CONCAT(DISTINCT categories.title) AS kategori, 
	GROUP_CONCAT(DISTINCT tags.name) AS tag, 
	blogs.thumbnail, 
	blogs.title, 
	blogs.body, 
	blogs.created_at, 
	blogs.updated_at 
	FROM blogs 
	LEFT JOIN users ON blogs.user_id = users.id 
	LEFT JOIN blog_categories ON blogs.id = blog_categories.id 
	LEFT JOIN categories ON blog_categories.id_category = categories.id 
	LEFT join blog_tags ON blogs.id = blog_tags.id_tag 
	LEFT JOIN tags ON blogs.id = tags.id 
	GROUP BY blogs.id
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
		var tag string
		if err := rows.Scan(&blog.Id, &blog.Username, &blog.Category, &tag, &blog.Title, &blog.Thumbnail, &blog.Body, &blog.CreatedAt, &blog.UpdatedAt); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		blog.TagName = SplitTags(tag)
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
	q := "SELECT blogs.id, users.username, GROUP_CONCAT(DISTINCT categories.title) AS kategori, GROUP_CONCAT(DISTINCT tags.name) AS tag, blogs.thumbnail, blogs.title, blogs.body, blogs.created_at, blogs.updated_at FROM blogs LEFT JOIN users ON blogs.user_id = users.id LEFT JOIN blog_categories ON blogs.id = blog_categories.id LEFT JOIN categories ON blog_categories.id_category = categories.id LEFT join blog_tags ON blogs.id = blog_tags.id_tag LEFT JOIN tags ON blogs.id = tags.id WHERE blogs.id = ?"
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

// TODO: refactor function
// create new blog - POST
func CreateBlog(c *fiber.Ctx) error {
	blog := new(blog.RequestCreateBlog)
	if err := c.BodyParser(blog); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	db := db.Connection()
	defer db.Close()

	q := "INSERT INTO `blogs` (`user_id`, `category_id`, `tag_id`, `title`, `thumbnail`, `body`) VALUES (?, ?, ?, ?, ?, ?)"

	stmt, err := db.Prepare(q)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
	defer stmt.Close()

	for _, tag := range blog.TagId {
		_, err := stmt.Exec(blog.UserId, blog.CategoryId, tag, blog.Title, blog.Thumbnail, blog.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data":    blog,
	})
}

// TODO: refactor function
func UpdateBlog(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := c.Params("userId")
	blog := new(blog.RequestUpdateBlog)
	if err := c.BodyParser(blog); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	db := db.Connection()
	defer db.Close()

	q := "UPDATE `blogs` SET `category_id` = ?,`tag_id` = ?,`title` = ?,`thumbnail` = ?,`body` = ?,`updated_at`= CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?"

	stmt, err := db.Prepare(q)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server errors",
		})
	}
	defer stmt.Close()

	// iterate for update tag
	for _, tag := range blog.TagId {
		_, err := stmt.Exec(blog.CategoryId, tag, blog.Title, blog.Thumbnail, blog.Body, id, userId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "succes update blog",
		"data":    blog,
	})
}
