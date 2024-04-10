package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/RizkiRdm/go-blog/pkg/models/users"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret")

// hash password
func GenerateHashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// compare password
func ComparePassword(password string, hash []byte) error {
	// Bandingkan password mentah dengan hash password.
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return err // Kembalikan error jika tidak cocok.
	}
	return nil // Password cocok, tidak ada error.
}

// generate jwt token for email
func GenerateToken(userId uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &users.Claims{
		UsernameId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// verify token
func VerifyTokenJWT(tokenString string) (*users.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &users.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*users.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// extract user_id from jwt token
func ExtractUserId(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}
	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in token")
	}
	return int(userId), nil
}
