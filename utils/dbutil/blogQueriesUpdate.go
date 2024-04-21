package dbutil

import (
	"database/sql"
	"mime/multipart"
)

// UPDATE BLOG FUNCTION
func UpdateBlog(tx *sql.Tx, form *multipart.Form, id string, userId int, thumbnailPath string) error {
	slug := SlugBlog(form.Value["title"][0])
	blogQuery := "UPDATE `blogs` SET `title`= ?, `slug`= ? `thumbnail` = ?, `body` = ? WHERE `id` = ? AND `user_id` = ?"
	result, err := tx.Exec(blogQuery, form.Value["title"][0], slug, thumbnailPath, form.Value["body"][0], id, userId)
	if err != nil {
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

	// update categoryId to blog_categories
	categoryQuery := "UPDATE `blog_categories` SET `id_blog`=?,`id_category`=? WHERE blogs.id = ? AND blogs.user_id = ?"
	if _, err := tx.Exec(categoryQuery, blogId, categoryId, id, userId); err != nil {
		return err
	}
	return nil
}

// GET CATEGORY ID
func getCategoryId(tx *sql.Tx, categoryName string) (int, error) {
	var categoryId int
	readCategoryQuery := "SELECT id FROM `categories` WHERE title = ?"
	if err := tx.QueryRow(readCategoryQuery, categoryName).Scan(&categoryId); err != nil {
		if err == sql.ErrNoRows {
			insertCategoryQuery := "INSERT INTO `categories` (title) VALUES (?)"
			if _, err := tx.Exec(insertCategoryQuery, categoryName); err != nil {
				return 0, err
			}
			// get category id
			if err := tx.QueryRow(readCategoryQuery, categoryName).Scan(&categoryId); err != nil {
				return 0, err
			}
			return categoryId, nil
		}
		return 0, err
	}
	return categoryId, nil
}

// GET OLD THUMBNAIL
func GetOldThumbnail(tx *sql.Tx, id string, userid int) (string, error) {
	var oldThumbnail string
	readThumbnailQuery := "SELECT `thumbnail` FROM `blogs` WHERE id = ? and user_id = ?"
	if err := tx.QueryRow(readThumbnailQuery, oldThumbnail).Scan(id, userid); err != nil {
		return "", err
	}
	return oldThumbnail, nil
}
