package users

import (
	"time"

	"github.com/RizkiRdm/go-blog/db"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Claims struct {
	UsernameId uint `json:"user_id"`
	jwt.RegisteredClaims
}

// function get user by email
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.Connection().QueryRow("SELECT `id`, `name`, `username`, `email`, `password` FROM `users` WHERE `email` = ?", email).Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password)
	return &user, err
}
