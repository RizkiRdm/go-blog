package middleware

import (
	"net/http"
	"strings"

	"github.com/RizkiRdm/go-blog/utils"
	"github.com/gofiber/fiber/v2"
)

func Midleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing authorization header ",
		})
	}

	tokenParts := strings.Split(authHeader, "")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token format",
		})
	}

	claims, err := utils.VerifyTokenJWT(tokenParts[1])
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message":    "invalid token",
			"messageErr": err.Error(),
		})
	}
	c.Locals("user", claims.UsernameId)
	return c.Next()
}
