package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("replace-with-secret")

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func JWTMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}
	var tokenStr string
	// expect "Bearer <token>"
	_, err := fmt.Sscanf(auth, "Bearer %s", &tokenStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid auth header")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// sub is float64 by default
		if sub, ok := claims["sub"].(float64); ok {
			c.Locals("user_id", uint(sub))
			return c.Next()
		}
	}
	return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
}
