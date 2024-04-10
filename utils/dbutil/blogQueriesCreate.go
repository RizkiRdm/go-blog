package dbutil

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
)

func CreateBlog(tx *sql.Tx, form *multipart.Form, userId int, thumbnailPath string) error {
	blogQuery := "INSERT INTO blogs (user_id, title, thumbnail, body) VALUES (?, ?, ?, ?)"
	result, err := tx.Exec(blogQuery, userId, form.Value["title"][0], thumbnailPath, form.Value["body"][0])
	if err != nil {
		// remove image
		if err := os.Remove(thumbnailPath); err != nil {
			return fmt.Errorf("failed to rollaback image")
		}
		return err
	}

	blogId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	categoryName := form.Value["category_name"][0]
	categoryId, err := getCategoryId(tx, categoryName)
	if err != nil {
		return err
	}

	// insert categoryId to blog_categories
	categoryQuery := "INSERT INTO `blog_categories`(`id_blog`, `id_category`) VALUES (?,?)"
	if _, err := tx.Exec(categoryQuery, blogId, categoryId); err != nil {
		return err
	}
	return nil
}
