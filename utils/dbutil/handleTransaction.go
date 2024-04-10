package dbutil

import (
	"database/sql"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func HandleTransactionError(tx *sql.Tx, c *fiber.Ctx, err error) error {
	tx.Rollback()
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"message": "transaction failed",
	})
}
