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

// NEW CODE

// small function
// add new blog
func addBlog(db *sql.Tx, userId int, thumbnail, title, body string) (int64, error) {
	q := "INSERT INTO blogs (user_id, title, thumbnail, body) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(q)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(userId, thumbnail, title, body)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// add tag blog
func addTag(db *sql.Tx, name string) (int64, error) {
	q := "INSERT INTO tags (name) VALUES (?) ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)"
	stmt, err := db.Prepare(q)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// add new category blog
func addCategory(db *sql.Tx, title string) (int64, error) {
	q := "INSERT INTO categories (title) VALUES (?) ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)"
	stmt, err := db.Prepare(q)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(title)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// associate blog with tag
func associateBlogWithTag(db *sql.Tx, blogId, tagId int) error {
	q := "INSERT INTO blog_tags (id_blog, id_tag) VALUES (?, ?) ON DUPLICATE KEY UPDATE id_blog = id_blog"
	_, err := db.Exec(q, blogId, tagId)
	return err
}

// associate blog with tag
func associateBlogWithCategory(db *sql.Tx, blogId, categoryId int) error {
	q := "INSERT INTO blog_categories (id, id_tag) VALUES (?, ?) ON DUPLICATE KEY UPDATE id = id"
	_, err := db.Exec(q, blogId, categoryId)
	return err
}

// handle insert data blog with tag & category
func HandleBlogWithDetails(c *fiber.Ctx) error {
	var req blog.RequestCreateBlog
	// handle body parse
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "error parsing request",
		})
	}

	// define variable db
	db := db.Connection()
	defer db.Close()

	// use transaction for make sure all operation success
	tx, err := db.Begin()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "internal server error",
		})
	}
	defer tx.Rollback()

	// add new blog
	blogId, err := addBlog(tx, req.UsernameId, req.Thumbnail, req.Title, req.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error add new blog",
		})
	}

	// add a tag and connect to blog
	for _, tag := range req.Tags {
		tagId, err := addTag(tx, tag)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "error adding tag",
			})
		}

		// associate tag and blog
		if err := associateBlogWithTag(tx, int(blogId), int(tagId)); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "error associating tage with blog",
			})
		}
	}

	// add category and connect to blog
	categoryId, err := addCategory(tx, req.Category)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error add category",
		})
	}

	if err := associateBlogWithCategory(tx, int(blogId), int(categoryId)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error associating cateogory to blog",
		})
	}

	// commit transaction if success
	if err := tx.Commit(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "error comitting transaction",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success",
		"data": fiber.Map{
			"username":  req.UsernameId,
			"thumbnail": req.Thumbnail,
			"title":     req.Title,
			"body":      req.Body,
			"category":  req.Category,
			"tags":      req.Tags,
		},
	})
}
