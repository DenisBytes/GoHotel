package api

import (
	"github.com/DenisBytes/GoHotel/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok{
		return ErrUnauthorized()
	}
	if user.IsAdmin==false{
		return ErrUnauthorized()
	}
	return c.Next()
}