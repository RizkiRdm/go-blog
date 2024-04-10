package users

import (
	"net/http"
	"time"

	"github.com/RizkiRdm/go-blog/db"
	"github.com/RizkiRdm/go-blog/pkg/models/users"
	"github.com/RizkiRdm/go-blog/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// register user - POST
func RegisterUser(c *fiber.Ctx) error {
	user := new(users.UsersCreateRequest)
	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "error parsing request user",
			"messageErr": err.Error(),
		})
	}

	// Hash password
	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed to hash password",
			"messageErr": err.Error(),
		})
	}

	// Define database connection
	db := db.Connection()
	defer db.Close()

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Check if email is already registered
	var emailCount int
	checkEmailQuery := "SELECT COUNT(*) FROM users WHERE email = ?"
	if err := tx.QueryRow(checkEmailQuery, user.Email).Scan(&emailCount); err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed to check email",
			"messageErr": err.Error(),
		})
	}

	if emailCount > 0 {
		tx.Rollback()
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "email already exists",
		})
	}

	// Insert user data into the database
	insertUserQuery := "INSERT INTO `users` (`name`, `username`, `email`, `password`) VALUES (?, ?, ?, ?)"
	if _, err := tx.Exec(insertUserQuery, user.Name, user.Username, user.Email, hashedPassword); err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed to insert data",
			"messageErr": err.Error(),
		})
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed to commit transaction",
			"messageErr": err.Error(),
		})
	}

	// Return success message
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "success create data",
		"data": fiber.Map{
			"name":     user.Name,
			"username": user.Username,
			"email":    user.Email,
			"password": hashedPassword,
		},
	})
}

// login user - POST
func LoginUser(c *fiber.Ctx) error {
	user := new(users.UserLoginRequest)
	if err := c.BodyParser(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "error parse body",
			"messageErr": err.Error(),
		})
	}

	// get user by email
	storedUser, err := users.GetUserByEmail(user.Email)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message":    "invalid credentials email",
			"messageErr": err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message":    "invalid credentials password",
			"messageErr": err.Error(),
		})
	}

	// generate token
	token, err := utils.GenerateToken(storedUser.Id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message":    "failed generate token",
			"messageErr": err.Error(),
		})
	}

	// create cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

// logout user - POST
func LogoutUser(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success logout",
	})
}
