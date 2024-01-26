package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DenisBytes/GoHotel/db"
	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler{
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email  string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User *types.User `json:"user"`
	Token string `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

// A handler should only do:
// - serialization of the incoming reuest (JSON)
// - DO some data fetching from db
// - call some business logic 
// - return the data back to the user

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error{
	var params AuthParams
	if err := c.BodyParser(&params); err !=nil{
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err!=nil{
		if errors.Is(err, mongo.ErrNoDocuments){
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.EncryptedPassword, params.Password){
		return invalidCredentials(c)
	}	

	token := CreateTokenFromUser(user)

	resp := AuthResponse{
		User: user,
		Token: token,
	}

	//TODO: when using thunder/
	// Set the X-API-TOKEN header in the response
	c.Set("X-API-TOKEN", token)

	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string{
	exp := time.Now().Add(time.Hour *4).Unix()	
	claims := jwt.MapClaims{
		"id": user.ID,
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err!=nil{
		fmt.Println("Failed to sign token with secret")
	}

	return tokenStr
}