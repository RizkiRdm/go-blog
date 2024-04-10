package dbutil

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadImage(c *fiber.Ctx, filename string) (string, error) {
	file, err := c.FormFile("thumbnail")
	if err != nil {
		return "", err
	}

	imageName := uuid.New().String() + filepath.Ext(file.Filename)
	thumbnailPath := filepath.Join("./uploads", imageName)

	// save image
	if err := c.SaveFile(file, thumbnailPath); err != nil {
		return "", err
	}

	return thumbnailPath, nil
}
