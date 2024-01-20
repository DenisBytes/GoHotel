package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error{
	fmt.Println("--- JWTAuth")
	token := c.Get("X-API-TOKEN")
	if len(token) == 0 {
		return fmt.Errorf("unauthorized")
	}
	if err:= parseJWTToken(token); err!=nil{
		return err
	}
	fmt.Println(token)
	return nil 
}

func parseJWTToken(tokenString string) error{
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
		return fmt.Errorf("Unauthorized")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
	}
	return fmt.Errorf("unauthorized")

}