package api

import (
	"fmt"
	"os"
	"time"

	"github.com/DenisBytes/GoHotel/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Decorator pattern : wrapping functions/data types
func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-API-TOKEN")
		if len(token) == 0 {
			return ErrUnauthorized()
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expFloat, ok := claims["exp"].(float64)
		if !ok {
			return fmt.Errorf("invalid expiration time")
		}

		exp := time.Unix(int64(expFloat), 0)
		if time.Now().After(exp) {
			return NewError(fiber.StatusUnauthorized, "token expired")
		}
		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnauthorized()
		}

		//TODO: this does not work with thunder/postman.
		//I have to put in the header manually the token eveytime.
		// Maybe in production works?
		// Set the current authenticated user to the context
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invali signing method", token.Header["alg"])
			return nil, ErrUnauthorized()
		}

		secret := os.Getenv("JWT_SECRET")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, ErrUnauthorized()
	}

	if !token.Valid {
		fmt.Println("Invalid token:", err)

		return nil, ErrUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrUnauthorized()
	}
	return claims, nil
}
