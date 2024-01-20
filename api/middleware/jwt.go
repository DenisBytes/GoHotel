package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error{
	token := c.Get("X-API-TOKEN")
	if len(token) == 0 {
		return fmt.Errorf("unauthorized")
	}
	claims, err:= validateToken(token)
	if err!=nil{
		return err
	}

	expFloat, ok := claims["exp"].(float64)
    if !ok {
        return fmt.Errorf("invalid expiration time")
    }

    exp := time.Unix(int64(expFloat), 0)
    if time.Now().After(exp){
		return fmt.Errorf("Token Expired")
	}
	return c.Next() 
}

func validateToken(tokenStr string) (jwt.MapClaims, error){
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invali signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
	
		secret := os.Getenv("JWT_SECRET")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, fmt.Errorf("Unauthorized")
	}

	if !token.Valid {
		fmt.Println("Invalid token:", err)

		return nil, fmt.Errorf("Unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}