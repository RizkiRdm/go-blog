package users

import (
	"net/http"

	"github.com/RizkiRdm/go-blog/db"
	"github.com/RizkiRdm/go-blog/pkg/models/users"
	"github.com/RizkiRdm/go-blog/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	user := new(users.UsersCreateRequest)
	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "error parsing request user",
			"messageErr": err.Error(),
		})
	}
	// hash password
	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed to hash password",
			"messageErr": err.Error(),
		})
	}
	// define variable db
	db := db.Connection()
	defer db.Close()

	// define transaction
	tx, err := db.Begin()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// check email are registered
	var emailCount int
	checkEmailQuery := "SELECT COUNT(*) FROM users WHERE email = ?"
	if err := tx.QueryRow(checkEmailQuery, user.Email).Scan(&emailCount); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed check email",
		})
	}

	if emailCount > 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "email already exist",
		})
	}

	// save data user to database
	insertUserQuery := "INSERT INTO `users`(`name`, `username`, `email`, `password`) VALUES ('?','?','?','?')"

	if _, err := tx.Exec(insertUserQuery, user.Name, user.Username, user.Email, hashedPassword); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed insert data",
			"messageErr": err.Error(),
		})
	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to commit transaction",
		})
	}

	// return success create data
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success create data",
		"data":    user,
	})
}
