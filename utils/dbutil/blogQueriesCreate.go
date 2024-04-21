package dbutil

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func CreateBlog(tx *sql.Tx, form *multipart.Form, userId int, thumbnailPath string) error {
	slug := SlugBlog(form.Value["title"][0])
	blogQuery := "INSERT INTO blogs (user_id, title, slug, thumbnail, body) VALUES (?, ?, ?, ?)"
	result, err := tx.Exec(blogQuery, userId, form.Value["title"][0], slug, thumbnailPath, form.Value["body"][0])
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

// create slug for insert and update blog
func SlugBlog(title string) string {
	slug := strings.ToLower(title)

	// delete character non-alphanumeric
	reg := regexp.MustCompile("[^a-z0-9]")
	slug = reg.ReplaceAllLiteralString(slug, "-")

	// remove hyphens at the beginning and at the end
	slug = strings.Trim(slug, "-")

	// add hypens at the end slug
	uuidString := uuid.New().String()
	slugWithUuid := fmt.Sprintf("%s-%s", slug, uuidString)

	return slugWithUuid
}
