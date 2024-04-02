package users

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
